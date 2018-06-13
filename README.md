# Baylands Geth Research Playground
## BinTrie Implementation
### Notes
In this implementation, we changed the radix of the state trie to 2 from 16. In
order to do this, we also needed to change the hex-prefix encoding.

Originally, the prefix only used the first two bits to specify if a node was a leaf
or not, but in order to avoid leaf keys that could not be stored in a byte (i.e.
any leaf node of length not divisible by 4), we used the first two bits to represent
an amount of padding, which was added to allow these leaf paths to be stored without
losing 1-3 bits of path information.
