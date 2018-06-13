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
	// "fmt"

	"bytes"
	"testing"
)

func TestDecodeBits(t *testing.T) {
	tests := []struct{ bin, bytes []byte}{
		{bin: []byte{}, bytes: []byte{}}, // empty bin means empty bytes
		{bin: []byte{1,1,1,1,1,1,1,1}, bytes: []byte{255}},
		{bin: []byte{0,0,0,0,0,0,0,1}, bytes: []byte{1}},
		{bin: []byte{1,1,1,1,1,1,1,1,0,0,0,0,0,0,0,1}, bytes: []byte{255, 1}},
		{bin: []byte{1,1,1,1,1,1,1,0}, bytes: []byte{254}},
	}
	for _, test := range tests {
		if b := testDecodeBits(test.bin); !bytes.Equal(b, test.bytes) {
			t.Errorf("testDecodeBits(%x) -> %x, want %x", test.bin, b, test.bytes)
		}
	}
}

func TestBinCompact(t *testing.T) {
    tests := []struct{ bin, compact []byte }{
        // empty keys, with and without terminator.
        {bin: []byte{}, compact: []byte{0x00}},
        {bin: []byte{2}, compact: []byte{0x20}},
        // odd length, no terminator
        {
            // hex: []byte{1, 2, 3, 4, 5},
            bin: []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1},
            compact: []byte{0x11, 0x23, 0x45},
        },
        // even length, no terminator
        {
            // hex: []byte{0, 1, 2, 3, 4, 5},
            bin: []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1},
            compact: []byte{0x00, 0x01, 0x23, 0x45},
        },
        // odd length, terminator
        {
            // hex: []byte{15, 1, 12, 11, 8, 16 /*term*/},
            bin: []byte{1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 2 /*term*/},
            compact: []byte{0x3f, 0x1c, 0xb8},
        },
        // even length, terminator
        {
            // hex: []byte{0, 15, 1, 12, 11, 8, 16 /*term*/},
            bin: []byte{0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 2 /*term*/},
            compact: []byte{0x20, 0x0f, 0x1c, 0xb8},
        },
				{
					bin: []byte{0, 1, 1},
					compact: []byte{0x56},
				},
				{
					bin: []byte{0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1, 2},
					compact: []byte{0x20, 0x64, 0x6f, 0x67, 0x65},
				},
				{
					bin: []byte{0, 1, 1, 0, 1},
					compact: []byte{0xc0, 0x68},
				},
				{
					bin: []byte{0, 1, 1, 0, 1, 2},
					compact: []byte{0xe0, 0x68},
				},
				{
					// needs 1 bit of padding, no terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1},
					compact: []byte{0x56, 0xca},
				},
				{
					// needs 1 bit of padding, with terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1, 2},
					compact: []byte{0x76, 0xca},
				},
				{
					// needs 2 bits of padding, no terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 1, 0},
					compact: []byte{0x96, 0xc8},
				},
				{
					// needs 2 bits of padding, with terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 2},
					compact: []byte{0xb6, 0xc8},
				},
				{
					// needs 3 bits of padding, no terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 0},
					compact: []byte{0xd6, 0xc0},
				},
				{
					// needs 3 bits of padding, with terminator
					bin: []byte{0, 1, 1, 0, 1, 1, 0, 0, 0, 2},
					compact: []byte{0xf6, 0xc0},
				},
    }
    for _, test := range tests {
				// fmt.Printf("compact: %x\n", test.compact)
				// fmt.Printf("bin: %+v\n", test.bin)
        if c := binToCompact(test.bin); !bytes.Equal(c, test.compact) {
            t.Errorf("binToCompact(%+v) -> %x, want %x", test.bin, c, test.compact)
        }
				// fmt.Printf("bin2: %+v\n", test.bin) // test.bin changed by here
        if h := compactToBin(test.compact); !bytes.Equal(h, test.bin) {
         		t.Errorf("compactToBin(%x) -> %+v, want %+v", test.compact, h, test.bin)
        }
				// fmt.Printf("bin3: %+v\n", test.bin)
    }
}

func TestBinKeybytes(t *testing.T) {
	tests := []struct{ key, binIn, binOut []byte }{
		{key: []byte{}, binIn: []byte{2}, binOut: []byte{2}}, // ???
		{key: []byte{}, binIn: []byte{}, binOut: []byte{2}}, // ???
		{
			key:    []byte{0x12, 0x34, 0x56},
			//hexIn:  []byte{1, 2, 3, 4, 5, 6, 16},
			binIn:  []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 2},
			//hexOut: []byte{1, 2, 3, 4, 5, 6, 16},
			binOut: []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 2},
		},
		{
			key:    []byte{0x12, 0x34, 0x5},
			//hexIn:  []byte{1, 2, 3, 4, 0, 5, 16},
			binIn:  []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 2},
			//hexOut: []byte{1, 2, 3, 4, 0, 5, 16},
			binOut: []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 2},
		},
		{
			key:    []byte{0x12, 0x34, 0x56},
			//hexIn:  []byte{1, 2, 3, 4, 5, 6},
			binIn:  []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0},
			//hexOut: []byte{1, 2, 3, 4, 5, 6, 16},
			binOut: []byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 2},
		},
	}
	for _, test := range tests {
		if h := keybytesToBin(test.key); !bytes.Equal(h, test.binOut) {
			t.Errorf("keybytesToBin(%x) -> %x, want %x", test.key, h, test.binOut)
		}
		if k := binToKeybytes(test.binIn); !bytes.Equal(k, test.key) {
			t.Errorf("binToKeybytes(%x) -> %x, want %x", test.binIn, k, test.key)
		}
	}
}

func BenchmarkBinToCompact(b *testing.B) {
	for t := 0; t < 10000; t++ {
		testBytes := []byte{0, 1, 1, 0, 1, 0, 2 /*term*/}
		for i := 0; i < b.N; i++ {
			binToCompact(testBytes)
		}
	}
}

func BenchmarkCompactToBin(b *testing.B) {
	for t := 0; t < 10000; t++ {
		testBytes := []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}
		for i := 0; i < b.N; i++ {
			compactToBin(testBytes)
		}
	}
}

func BenchmarkKeybytesToBin(b *testing.B) {
	for t := 0; t < 10000; t++ {
		testBytes := []byte{7, 6, 6, 5, 7, 2, 6, 2, 16}
		for i := 0; i < b.N; i++ {
			keybytesToBin(testBytes)
		}
	}
}

func BenchmarkBinToKeybytes(b *testing.B) {
	for t := 0; t < 10000; t++ {
		testBytes := []byte{0, 1, 1, 1, 1, 0, 0, 0, 2}
		for i := 0; i < b.N; i++ {
			binToKeybytes(testBytes)
		}
	}
}
