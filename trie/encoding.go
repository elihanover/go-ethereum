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

// import "fmt"

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
	if hasTerm(hex) {
		terminator = 1
		hex = hex[:len(hex)-1]
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

func binToCompactOG(b []byte) []byte {
	bin := make([]byte, len(b))
	copy(bin, b)
	terminator := byte(0)
	if hasTerm(bin) {
		terminator = 1
		bin = bin[:len(bin)-1] // take off terminator
	}
	// if len(bin) not divisible by 4, pad the ending with 0's,
	// and tell first 2 bits of prefix how much we padded
	padding := (4 - len(bin) % 4) % 4
	// add padding
	if (len(bin)%4) != 0 {
		for i := 0; i < padding; i++ {
			bin = append(bin, 0x0)
		}
	} // doesn't work for len(bin) < 4 and len(bin) % 4 != 0 ????
	buf := make([]byte, len(bin)/8+1)
	buf[0] = terminator << 5
	buf[0] |= byte(padding) << 6
	odd := (len(bin)/4&1) == 1
	if odd { // if odd
		// fmt.Println("is odd")
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= bin[0] << 3 // first 4 bits is contained in the first byte
		buf[0] |= bin[1] << 2// but what if we only have one bit?
		buf[0] |= bin[2] << 1
		buf[0] |= bin[3]
		bin = bin[4:]
	}
	// fmt.Printf("buf[0] 2: %x\n", buf[0])
	decodeBits(bin, buf[1:])
	// fmt.Printf("binside: %+v\n", bin)
	return buf
}

func binToCompact(b []byte) []byte {
	terminator := byte(0)
	if hasTerm(b) {
		terminator = 1
		b = b[:len(b)-1] // take off terminator
	}
	// if len(bin) not divisible by 4, pad the ending with 0's,
	// and tell first 2 bits of prefix how much we padded.
	padding := (4 - len(b) % 4) % 4
	bin := make([]byte, len(b)+padding)
	copy(bin, b)
	// doesn't work for len(bin) < 4 and len(bin) % 4 != 0
	buf := make([]byte, len(bin)/8+1)
	buf[0] = terminator << 5
	buf[0] |= byte(padding) << 6
	// fmt.Printf("len bin: %d\n", len(bin))
	// odd := (len(bin)/4&1) == 1
	if (len(bin)/4&1) == 1 { // if odd
		// fmt.Println("is odd")
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= bin[0] << 3 // first 4 bits is contained in the first byte
		buf[0] |= bin[1] << 2// but what if we only have one bit?
		buf[0] |= bin[2] << 1
		buf[0] |= bin[3]
		bin = bin[4:]
	}
	// fmt.Printf("buf[0] 2: %x\n", buf[0])
	decodeBits(bin, buf[1:])
	// fmt.Printf("binside: %+v\n", bin)
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
 	// fmt.Printf("base: %+v\n", base)
	base = base[:len(base)-1] // take off terminator bit
	// apply odd flag
	chop := 4 * (2 - int(base[3]))
	// check for terminator
	terminator := base[2] != 0
	// check for end padding
	pad := int((base[0] << 1) + base[1])
	// fmt.Printf("pad: %d\n", pad)
	base = base[chop:len(base)-pad]
	// fmt.Printf("base before terminator stuff: %+v\n", base)
	// apply terminator flag
	if terminator {
		base = append(base, 2) // terminator
	}
	// fmt.Printf("base: %+v\n", base)
	return base
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

// binToKeybytes turns binary encoded bytes into key bytes.
// This can only be used for keys of length % 8 == 0.
func binToKeybytes(bin []byte) []byte {
	// fmt.Printf("bin2kyb: %+v\n", bin)
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
	// fmt.Printf("decode: %+v\n", bits)
	// fmt.Printf("into %d bytes\n", len(bytes))
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
