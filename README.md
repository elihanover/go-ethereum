# Go-Ethereum Binary State Trie Implementation
## TL;DR
By using a binary Merkle Patricia tree and removing redundant hash nodes from Merkle proofs, we can reduce the size of proofs by over 40%.

## Context
As mentioned [here](https://ethresear.ch/t/a-two-layer-account-trie-inside-a-single-layer-trie/210), reducing the radix of the state trie to 2 would decrease the size of the light client’s optimal merkle proof by a factor of ~3.75.

3.75 = 15/4, where we have 15 times less sibling nodes per layer and 4x more layers.  The cost of this is a trie with 87% more nodes (again, not accounting for leaf/extension nodes), an added burden for full nodes to store, construct merkle proofs, and hash from.  Note that the size of branch nodes is substantially (~5x) smaller using a binary trie, so 87% may not be as dire as it appears.

## What We Did
### "Hashless" Merkle Proofs
In its current implementation, Merkle proofs are constructed as key/value mappings where the hash of the node is the key to the nodes RLP encoding.  However, can avoid storing the hash keys of the database entirely (this approach is hinted at in the [Ethereum Wishlist](https://github.com/ethereum/wiki/wiki/Wishlist#trie)).  This is vital in a radix 2 implementation as the Merkle branch can be up to 4x longer that in the radix 16 implementation, meaning 4x the hash keys to store in the proof.  Even with significantly smaller nodes, this adds up.

#### Solution
Instead of using the proof database hash keys to check for the next valid proof node, we instead compare the hash of the next proof node with the next child hash node.  This allows us to bypass storing the hash keys for each proof node.

### Bin-Prefix "BP" Encoding
Besides changing some node properties, the one main change that was needed was a change to the encoding of the node paths.
A node's path can be encoded in either HP (referred to as *compact* encoding in the code) or hex (binary in our case).

#### Compact encoding looks as follows:

##### (Hex) Prefix: high nibble of first byte
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

    Hex: "43c1T" (note that each hex character is actually a byte after decoding, so T corresponds to a terminator BYTE)
    Compact: [0010 0000] [0100 0011] [1100 0001]

#### Problem:
In order to change the radix of the state trie to 2, we needed to change the HP encoding scheme slightly.  With hex characters, you need 2 characters per byte.  As seen in the traditional HP encoding, all you need to deal with the odd case is use some padding bits.

However in the case of binary paths, you need 8 characters per byte, which means you have to account for 7 cases in which the path does not fit evenly into a byte.

To account for this, we used the first two bits of the prefix to represent some extra padding that was added to the end of the path.  We took advantage of the four padding bits implemented by the HP encoding, and added 0-3 bits of padding.

##### (Bin) Prefix: high nibble of first byte
    bit 0: padding bit
    bit 1: padding bit
    bit 2: terminator bit (set to 1 if path ends in 16)
    bit 3: odd/even bit (set to 1 if path is odd number of hex chars)


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



## Results
##### Hex Trie: Random 500 Nodes
    Ave Proof Nodes: 4.3
    Ave Proof Size: 1378.8 bytes

##### Bin Trie: Random 500 Nodes
    Ave Proof Nodes: 11.6
    Ave Proof Size: 768.4 bytes (44% improvement)

##### Hex Trie: Random 5000 Nodes
    Ave Proof Nodes: 5.2
    Ave Proof Size: 1759 bytes

##### Bin Trie: Random 5000 Nodes
    Ave Proof Nodes: 15.4
    Ave Proof Size: 1014 bytes (42% improvement)

##### Hex Trie: Random 50000 Nodes
    Ave Proof Nodes: 5.4
    Ave Proof Size: 2182.6 bytes

##### Bin Trie: Random 50000 Nodes
    Ave Proof Nodes: 18.6
    Ave Proof Size: 1235.4 bytes (43% improvement)

##### Hex Trie: Random 500000 Nodes
    Ave Proof Nodes: 6.2
    Ave Proof Size: 2663 bytes

##### Bin Trie: Random 500000 Nodes
    Ave Proof Nodes: 21
    Ave Proof Size: 1406.8 bytes (47% improvement)
