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

// Trie keys are dealt with in three distinct encodings:
//
// KEYBYTES encoding contains the actual key and nothing else. This encoding is the
// input to most API functions.
//
// HEX encoding contains one byte for each nibble of the key and an optional trailing
// 'terminator' byte of value 0x10 which indicates whether or not the node at the key
// contains a value. Hex key encoding is used for nodes loaded in memory because it's
// convenient to access.
//
// COMPACT encoding is defined by the Ethereum Yellow Paper (it's called "hex prefix
// encoding" there) and contains the bytes of the key and a flag. The high nibble of the
// first byte contains the flag; the lowest bit encoding the oddness of the length and
// the second-lowest encoding whether the node at the key is a value node. The low nibble
// of the first byte is zero in the case of an even number of nibbles and the first nibble
// in the case of an odd number. All remaining nibbles (now an even number) fit properly
// into the remaining bytes. Compact encoding is used for nodes stored on disk.

// Original
func hexToCompact(hex []byte) []byte {
	terminator := byte(0)
	if hasTerm(hex) { // IF THIS IS A LEAF
		terminator = 1
		hex = hex[:len(hex)-1] // GET RID OF TERMINATOR, WHY DID WE NEED IT IN MEMORY?
	}
	buf := make([]byte, len(hex)/2+1)
	buf[0] = terminator << 5 // the flag byte
	if len(hex)&1 == 1 { // if odd length
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= hex[0] // first nibble is contained in the first byte
		hex = hex[1:]
	}
	decodeNibbles(hex, buf[1:])
	return buf
}

// Modified
func binToCompact(bin []byte) []byte {
	terminator := byte(0)
	if hasTerm(bin) {
		terminator = 1
		bin = bin[:len(bin)-1]
	}
	buf := make([]byte, len(bin)/8+1)
	buf[0] = terminator << 5 // terminator flag byte.
	if (len(bin)/4)&1 == 1 { // if odd
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= bin[0] << 3 // first 4 bits is contained in the first byte
		buf[0] |= bin[1] << 2// but what if we only have one bit?
		buf[0] |= bin[2] << 1
		buf[0] |= bin[3]
		bin = bin[4:]
	}
	decodeBits(bin, buf[1:])
	return buf
}

// original
func compactToHex(compact []byte) []byte {
	base := keybytesToHex(compact)
	base = base[:len(base)-1]
	// apply terminator flag
	if base[0] >= 2 {
		base = append(base, 16)
	}
	// apply odd flag
	chop := 2 - base[0]&1
	return base[chop:]
}

// modified
func compactToBin(compact []byte) []byte {
	base := keybytesToBin(compact)
	prefix := base[0] << 3 | base[1] << 2 | base[2] << 1 | base[3]
	base = base[:len(base)-1] // take off terminator bit
	// apply terminator flag
	if prefix >= 2 {
		base = append(base, 2) // terminator
	}
	// apply odd flag
	chop := 4 * (2 - prefix&1)
	return base[chop:]
}

// original
func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1 // EXTRA BYTE IS TERMINATOR
	var nibbles = make([]byte, l) // OH BECAUSE NIBBLE ARE HALF THE SIZE, BUT THE +1?
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}


// modified
func keybytesToBin(str []byte) []byte {
	l := len(str) * 8 + 1
	var bits = make([]byte, l)
	for i, b := range str {
		for j := 0; j < 8; j++ {
			bits[8*i+j] = (b >> uint(7-j)) & 0x1
		}
	}
	bits[l-1] = 2 // set terminator bit
	return bits
}

//ORIGINAL
// hexToKeybytes turns hex nibbles into key bytes.
// This can only be used for keys of even length.
func hexToKeybytes(hex []byte) []byte {
	if hasTerm(hex) {
		hex = hex[:len(hex)-1]
	}
	if len(hex)&1 != 0 {
		panic("can't convert hex key of odd length")
	}
	key := make([]byte, (len(hex)+1)/2)
	decodeNibbles(hex, key)
	return key
}

// MODIFIED
// binToKeybytes turns binary encoded bytes into key bytes.
// This can only be used for keys of even length.
func binToKeybytes(bin []byte) []byte {
	if hasTerm(bin) { // does this have terminator flag?
		bin = bin[:len(bin)-1] // if so, drop it
	}
	if len(bin)%8 != 0 {
		panic("can't convert bin key, not divisible by 8")
	}
	key := make([]byte, len(bin)/8)
	decodeBits(bin, key)
	return key
}

// decodeBits into one slice of bytes.
func decodeBits(bits []byte, bytes []byte) []byte {
	for by := 0; by < len(bytes); by++ {
		for bt := 0; bt < 8; bt++ { // decode next 8 bits per byte
			bytes[by] |= bits[8*by+7-bt] << uint(bt)
		}
	}
	return bytes
}

func testDecodeBits(bits []byte) []byte {
	bytes := make([]byte, len(bits)/8)
	for by := 0; by < len(bytes); by++ {
		for bt := 0; bt < 8; bt++ { // decode next 8 bits per byte
			bytes[by] |= bits[8*by+7-bt] << uint(bt)
		}
	}
	return bytes
}

// decodeNibbles decodes the nibbles into an array of bytes.
func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}

// prefixLen returns the length of the common prefix of a and b.
func prefixLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

// DON'T NEED (?)
// prefixBitLen returns the length of the common prefix of a and b in terms of BITS.
func prefixBitLen(a, b []byte) int {
	var lenA = len(a)*8
	var lenB = len(b)*8
	var i, length = 0, lenA // 8 bits
	if lenB < length {
		length = lenB
	}
	for ; i < length; i++ {
		// check bit by bit for difference
		var byteIndex = i/8
		if a[byteIndex] >> uint(7-i%8) != b[byteIndex] >> uint(7-i%8) { // if
			break
		}
	}
	return i
}

// hasTerm returns whether a hex key has the terminator flag.
func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == 2 // terminator 2 instead of 16
}
