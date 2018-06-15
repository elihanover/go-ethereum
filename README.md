# Binary Patricia State Trie
## Context
As mentioned [here](https://ethresear.ch/t/a-two-layer-account-trie-inside-a-single-layer-trie/210), reducing the radix of the state trie to 2 would decrease the size of the light client’s merkle proof by a factor of ~3.75.

3.75 = 15/4, where we have 15 times less sibling nodes per layer and 4x more layers. although extension and leaf nodes complicate this, as shown in our tests.  This results in a trie with ~1.87x more nodes (again, not accounting for leaf/extension nodes), an added burden for full nodes to store, construct merkle proofs, and hash from.  So in the end we're left with a tradeoff and a consequent discussion of how to handle it.

## What We Did
Besides changing some node properties, the one main change that was needed was a change to the encoding of the node paths.

A node's path can be encoded in either HP (referred to as *compact* encoding in the code) or hex (binary in our case).

#### Compact encoding looks as follows:

##### Prefix: high nibble of first byte
    bit 0: no meaning
    bit 1: no meaning
    bit 2: terminator bit (set to 1 if path ends in 16)
    bit 3: odd/even bit (set to 1 if path is odd number of hex chars)

##### (Optional) Padding: low nibble of first byte
    If odd path length, pad this nibble with 0's in order to fit path evenly into bytes.

    Otherwise, set this nibble to the first hex character in the path.

##### Path:
    The next n bytes represent the next 2n hex characters of the path.

###### Examples:
    Hex: "a56cc"
    Compact: [0001 1010] [0101 0110] [1100 1100]

    Hex: "43c1f"
    Compact: [0010 0000] [0100 0011] [1100 0001]

#### Problem:
In order to change the radix of the state trie to 2, we needed to change the HP encoding scheme slightly.  With hex characters, you need 2 characters per byte.  As seen in the traditional HP encoding, all you need to deal with the odd case is use some padding bits.

However in the case of binary paths, you need 8 characters per byte, which means you have to account for 7 cases in which the path does not fit evenly into a byte.

To account for this, we used the first two bits of the prefix to represent some extra padding that was added to the end of the path.  We took advantage of the four padding bits implemented by the HP encoding, and added 0-3 bits of padding.

###### Example 1:
    Bin: 110010
    Compact': [1000 0000] [1100 1000]

  Note that the last two bits are the extra padding mentioned above, and are removed from the path when converted back to bin encoding.

  The first two bits of the prefix indicate that we have padded two bits onto the end of the path, and the last two bits of the path are treated the same as traditional HP encoding.

  **Note that the odd/even bit of the prefix assumes the padding has already been added.**


###### Example 2:
    Bin: 01111002
    Compact': [0110 0000] [0111 1000]

  In this case we've padded one bit onto the end of the path, have indicated there is a terminator bit, and that we need padding after the prefix to fit evenly into bytes.

  **Note that we changed the terminator from 16 to 2. This is the access the child node at index 2 of a branch node.**
  

## Benchmarks

##### Radix 16: Encoding Benchmark Results:
    BenchmarkHexToCompact-4    	    50000000	        30.0 ns/op
    BenchmarkCompactToHex-4    	    30000000	        37.6 ns/op
    BenchmarkKeybytesToHex-4   	    30000000	        50.4 ns/op
    BenchmarkHexToKeybytes-4   	    50000000	        25.6 ns/op

##### Radix 2: Encoding Benchmark Results:
    BenchmarkBinToCompact-4        	50000000	        27.4 ns/op
    BenchmarkCompactToBin-4        	20000000	        76.3 ns/op
    BenchmarkKeybytesToBin-4       	20000000	        90.6 ns/op
    BenchmarkBinToKeybytes-4       	50000000	        24.0 ns/op

##### Radix 16: Proof Benchmark Results:
    BenchmarkProve-4           	    2000	    801100 ns/op
    BenchmarkVerifyProof100-4     	  200000	     10789 ns/op
    BenchmarkVerifyProof1000-4     	  200000	     11673 ns/op
    BenchmarkVerifyProof10000-4    	  200000	     10828 ns/op
    BenchmarkVerifyProof100000-4   	  200000	     12626 ns/op
    
##### Radix 2: Proof Benchmark Results:
    BenchmarkProve-4               	     300	     3633869 ns/op
    BenchmarkVerifyProof100-4         100000	       10826 ns/op
    BenchmarkVerifyProof1000-4     	  100000	       11997 ns/op
    BenchmarkVerifyProof10000-4    	  200000	       10998 ns/op
    BenchmarkVerifyProof100000-4   	  100000	       14474 ns/op
    
##### Radix 16: Trie Benchmark Results:
    BenchmarkGet-4                 	 5000000	       247 ns/op
    BenchmarkGetDB-4               	10000000	       168 ns/op
    BenchmarkUpdateBE-4            	 1000000	      1863 ns/op
    BenchmarkUpdateLE-4            	 1000000	      2413 ns/op
    BenchmarkHash-4                	  300000	      4149 ns/op
    
##### Radix 2: Trie Benchmark Results:
    BenchmarkGet-4                 	 3000000	         519 ns/op
    BenchmarkGetDB-4               	 3000000	         406 ns/op
    BenchmarkUpdateBE-4            	 1000000	        2029 ns/op
    BenchmarkUpdateLE-4            	 1000000	        2707 ns/op
    BenchmarkHash-4                	  300000	        6757 ns/op

