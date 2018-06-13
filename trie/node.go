// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var indices = []string{"0", "1", "[3]"}

type node interface {
	fstring(string) string
	cache() (hashNode, bool)
	canUnload(cachegen, cachelimit uint16) bool
}

type (
	fullNode struct {
		Children [3]node // Actual trie node data to encode/decode (needs custom encoder)
		flags    nodeFlag
	}
	shortNode struct {
		Key   []byte
		Val   node
		flags nodeFlag // NEED TO CHANGE THESE?
	}
	hashNode  []byte
	valueNode []byte
)

// EncodeRLP encodes a full node into the consensus RLP format.
func (n *fullNode) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, n.Children)
}

func (n *fullNode) copy() *fullNode   { copy := *n; return &copy }
func (n *shortNode) copy() *shortNode { copy := *n; return &copy }

// BinKey returns a binary encoded shortNode key
func (n *shortNode) Bin0Key() []byte {
	// fmt.Printf("\n\n\n\nHELLO\n%x\n\n\n\n", n.Key)
	l := len(n.Key) * 4 + 1
	bin := make([]byte, l)
	for i := 0; i < len(n.Key)-1; i++ {
		for j := 0; j < 4; j++ {
			bin[4*i+j] = n.Key[i] & (0x8 >> uint(j))
		}
	}
	bin[l-1] = 2
	return bin
}

// nodeFlag contains caching-related metadata about a node.
type nodeFlag struct {
	hash  hashNode // cached hash of the node (may be nil)
	gen   uint16   // cache generation counter
	dirty bool     // whether the node has changes that must be written to the database
}

// canUnload tells whether a node can be unloaded.
func (n *nodeFlag) canUnload(cachegen, cachelimit uint16) bool {
	return !n.dirty && cachegen-n.gen >= cachelimit
}

func (n *fullNode) canUnload(gen, limit uint16) bool  { return n.flags.canUnload(gen, limit) }
func (n *shortNode) canUnload(gen, limit uint16) bool { return n.flags.canUnload(gen, limit) }
func (n hashNode) canUnload(uint16, uint16) bool      { return false }
func (n valueNode) canUnload(uint16, uint16) bool     { return false }

func (n *fullNode) cache() (hashNode, bool)  { return n.flags.hash, n.flags.dirty }
func (n *shortNode) cache() (hashNode, bool) { return n.flags.hash, n.flags.dirty }
func (n hashNode) cache() (hashNode, bool)   { return nil, true }
func (n valueNode) cache() (hashNode, bool)  { return nil, true }

// Pretty printing.
func (n *fullNode) String() string  { return n.fstring("") }
func (n *shortNode) String() string { return n.fstring("") }
func (n hashNode) String() string   { return n.fstring("") }
func (n valueNode) String() string  { return n.fstring("") }

func (n *fullNode) fstring(ind string) string {
	resp := fmt.Sprintf("[\n%s  ", ind)
	for i, node := range n.Children {
		if node == nil {
			resp += fmt.Sprintf("%s: <nil> ", indices[i])
		} else {
			resp += fmt.Sprintf("%s: %v", indices[i], node.fstring(ind+"  "))
		}
	}
	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *shortNode) fstring(ind string) string {
	return fmt.Sprintf("{%x: %v} ", n.Key, n.Val.fstring(ind+"  "))
}
func (n hashNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", []byte(n))
}
func (n valueNode) fstring(ind string) string {
	return fmt.Sprintf("%x ", []byte(n))
}

func mustDecodeNode(hash, buf []byte, cachegen uint16) node {
	n, err := decodeNode(hash, buf, cachegen)
	if err != nil {
		panic(fmt.Sprintf("node %x: %v", hash, err))
	}
	// fmt.Printf("Got: %+v", n)
	return n
}

// decodeNode parses the RLP encoding of a trie node.
func decodeNode(hash, buf []byte, cachegen uint16) (node, error) {
	if len(buf) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	elems, _, err := rlp.SplitList(buf)
	if err != nil {
		fmt.Println("Error in splitlist")
		return nil, fmt.Errorf("decode error: %v", err)
	}
	switch c, _ := rlp.CountValues(elems); c {
	case 2:
		// fmt.Printf("Node is short\n")
		n, err := decodeShort(hash, elems, cachegen)
		return n, wrapError(err, "short")
	case 3:
		// fmt.Printf("Node is full\n")
		n, err := decodeFull(hash, elems, cachegen)
		return n, wrapError(err, "full")
	default:
		// fmt.Printf("CASE4")
		return nil, fmt.Errorf("invalid number of list elements: %v", c)
	}
}

func decodeShort(hash, elems []byte, cachegen uint16) (node, error) {
	// fmt.Printf("Decode short: %+v", hash)
	kbuf, rest, err := rlp.SplitString(elems)
	// fmt.Printf("\nkbuf: %x, rest: %x\n", kbuf, rest)
	if err != nil {
		return nil, err
	}
	flag := nodeFlag{hash: hash, gen: cachegen}
	// fmt.Printf("kbuf: %x", kbuf)
	key := compactToBin(kbuf) // key returned is wrong, missing pieces on the end
	// fmt.Printf("\nshortnode bin key: %+v\n", key)
	if hasTerm(key) {
		// fmt.Printf("\nLEAFNODE\n")
		// value node
		val, _, err := rlp.SplitString(rest) // get value from leaf node
		if err != nil {
			return nil, fmt.Errorf("invalid value node: %v", err)
		}
		// fmt.Printf("\nDECODEVAL: %x\nString: %s\n", val, string(val))
		return &shortNode{key, append(valueNode{}, val...), flag}, nil
	}
	r, _, err := decodeRef(rest, cachegen) // get value from extension node
	if err != nil {
		return nil, wrapError(err, "val")
	}
	return &shortNode{key, r, flag}, nil
}

func decodeFull(hash, elems []byte, cachegen uint16) (*fullNode, error) {
	// fmt.Printf("\nDecode full: %x", hash)
	n := &fullNode{flags: nodeFlag{hash: hash, gen: cachegen}}
	for i := 0; i < 2; i++ {
		cld, rest, err := decodeRef(elems, cachegen)
		if err != nil {
			return n, wrapError(err, fmt.Sprintf("[%d]", i))
		}
		n.Children[i], elems = cld, rest
	}
	val, _, err := rlp.SplitString(elems)
	if err != nil {
		return n, err
	}
	if len(val) > 0 {
		n.Children[2] = append(valueNode{}, val...)
	}
	return n, nil
}

const hashLen = len(common.Hash{})

func decodeRef(buf []byte, cachegen uint16) (node, []byte, error) {
	kind, val, rest, err := rlp.Split(buf)
	// fmt.Printf("\nkind: %s\nval: %+v\nrest: %+v\n", kind, val, rest)
	if err != nil {
		// fmt.Println("T0")
		return nil, buf, err
	}
	switch {
	case kind == rlp.List:
		// fmt.Println("T1")
		// 'embedded' node reference. The encoding must be smaller
		// than a hash in order to be valid.
		if size := len(buf) - len(rest); size > hashLen {
			err := fmt.Errorf("oversized embedded node (size is %d bytes, want size < %d)", size, hashLen)
			return nil, buf, err
		}
		n, err := decodeNode(nil, buf, cachegen)
		return n, rest, err
	case kind == rlp.String && len(val) == 0:
		// fmt.Println("T2")
		// empty node
		return nil, rest, nil
	case kind == rlp.String && len(val) == 32:// what is this 32?  and should it be changed??
		// fmt.Println("T3")
		return append(hashNode{}, val...), rest, nil
	default:
		return nil, nil, fmt.Errorf("invalid RLP string size %d (want 0 or 32)", len(val))
	}
}

// wraps a decoding error with information about the path to the
// invalid child node (for debugging encoding issues).
type decodeError struct {
	what  error
	stack []string
}

func wrapError(err error, ctx string) error {
	if err == nil {
		return nil
	}
	if decErr, ok := err.(*decodeError); ok {
		decErr.stack = append(decErr.stack, ctx)
		return decErr
	}
	return &decodeError{err, []string{ctx}}
}

func (err *decodeError) Error() string {
	return fmt.Sprintf("%v (decode path: %s)", err.what, strings.Join(err.stack, "<-"))
}
