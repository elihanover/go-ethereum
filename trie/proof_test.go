// Copyright 2015 The go-ethereum Authors
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

	"bytes"
	crand "crypto/rand"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethdb"
)

func init() {
	mrand.Seed(time.Now().Unix())
}

func TestProof(t *testing.T) {
	trie, vals := randomTrie(500)
	root := trie.Hash()
	aveproofsize := 0
	ps := 0
	for _, kv := range vals {
		ps += 1

		proof, err := trie.Prove(kv.k, 0)
		if err != nil {
			t.Fatalf("missing key %x while constructing proof", kv.k)
		}

		proofsize := 0
		for _, p := range proof {
			proofsize += len(p)
		}
		aveproofsize += proofsize
		// fmt.Printf("proofsize: %+v\n", proofsize)

		val, err, _ := VerifyProof(root, kv.k, proof)
		if err != nil {
			t.Fatalf("VerifyProof error for key %x: %v\nraw proof: %v", kv.k, err, proof)
		}
		if !bytes.Equal(val, kv.v) {
			t.Fatalf("VerifyProof returned wrong value for key %x: got %x, want %x", kv.k, val, kv.v)
		}
	}
	fmt.Println(aveproofsize)
	fmt.Println(ps)
	fmt.Println(aveproofsize/ps)
	t.Fatalf("blah")
}

func TestOneElementProof(t *testing.T) {
	trie := new(Trie)
	updateString(trie, "k", "v")
	proof, err := trie.Prove([]byte("k"), 0)
	if len(proof) != 1 {
		t.Error("proof should have one element")
	}
	val, err, _ := VerifyProof(trie.Hash(), []byte("k"), proof)
	if err != nil {
		t.Fatalf("VerifyProof error: %v\nproof hashes: %v", err, proof)
	}
	if !bytes.Equal(val, []byte("v")) {
		t.Fatalf("VerifyProof returned wrong value: got %x, want 'k'", val)
	}
}

func TestVerifyBadProof(t *testing.T) {
	trie, vals := randomTrie(800)
	root := trie.Hash()
	for _, kv := range vals {
		proof, _ := trie.Prove(kv.k, 0)
		if len(proof) == 0 {
			t.Fatal("zero length proof")
		}

		node := proof[0]
		mutateByte(node)
		proof = proof[1:]
		proof = append(proof, node)
		if _, err, _ := VerifyProof(root, kv.k, proof); err == nil {
			t.Fatalf("expected proof to fail for key %x", kv.k)
		}
	}
}

// mutateByte changes one byte in b.
func mutateByte(b []byte) {
	for r := mrand.Intn(len(b)); ; {
		new := byte(mrand.Intn(255))
		if new != b[r] {
			b[r] = new
			break
		}
	}
}

// Get the distribution of nodes in a randomly generated trie
func BenchmarkAll(b *testing.B) {
	averagedDist := []int{0,0,0,0}
	averageProofsize := 0

	samples := 1000
	trieSize := 100

	for i := 0; i < samples; i++ {

		// generate random trie
		trie, vals := randomTrie(trieSize)

		// get node type distribution and add to average
		dist := trie.GetNodeTypeDistribution(trie.root)
		averagedDist[0] += dist[0]
		averagedDist[1] += dist[1]
		averagedDist[2] += dist[2]
		averagedDist[3] += dist[3]

		// get proof size
		var keys []string
		for k := range vals {
			keys = append(keys, k)
		}

		// Construct proof and get the size of it
		// b.ResetTimer()
		for i := 0; i < b.N; i++ {
			kv := vals[keys[i%len(keys)]]
			proof, _ := trie.Prove(kv.k, 0)
			if len(proof) == 0 {
				b.Fatalf("zero length proof for %x", kv.k)
			}

			averageProofsize += len(proof)
		}
	}

	// report back what we want
	fmt.Printf("\nFor tree of size %d:\n", samples)
	averagedDist[0] /= samples
	averagedDist[1] /= samples
	averagedDist[2] /= samples
	averagedDist[3] /= samples
	fmt.Printf("avedist: %+v\n", averagedDist)
	fmt.Printf("aveproof: %d\n", averageProofsize/samples)
}
//
// func BenchmarkProve(b *testing.B) {
// 	trie, vals := randomTrie(100)
// 	var keys []string
// 	for k := range vals {
// 		keys = append(keys, k)
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		kv := vals[keys[i%len(keys)]]
// 		proofs := ethdb.NewMemDatabase()
// 		if trie.Prove(kv.k, 0, proofs); len(proofs.Keys()) == 0 {
// 			b.Fatalf("zero length proof for %x", kv.k)
// 		}
// 	}
// }
//
// func BenchmarkVerifyProof(b *testing.B) {
// 	trie, vals := randomTrie(100)
// 	root := trie.Hash()
// 	var keys []string
// 	var proofs []*ethdb.MemDatabase
// 	for k := range vals {
// 		keys = append(keys, k)
// 		proof := ethdb.NewMemDatabase()
// 		trie.Prove([]byte(k), 0, proof)
// 		proofs = append(proofs, proof)
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		im := i % len(keys)
// 		if _, err, _ := VerifyProof(root, []byte(keys[im]), proofs[im]); err != nil {
// 			b.Fatalf("key %x: %v", keys[im], err)
// 		}
// 	}
// }
//
// func BenchmarkVerifyProof1000(b *testing.B) {
// 	trie, vals := randomTrie(1000)
// 	root := trie.Hash()
// 	var keys []string
// 	var proofs []*ethdb.MemDatabase
// 	for k := range vals {
// 		keys = append(keys, k)
// 		proof := ethdb.NewMemDatabase()
// 		trie.Prove([]byte(k), 0, proof)
// 		proofs = append(proofs, proof)
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		im := i % len(keys)
// 		if _, err, _ := VerifyProof(root, []byte(keys[im]), proofs[im]); err != nil {
// 			b.Fatalf("key %x: %v", keys[im], err)
// 		}
// 	}
// }
//
// func BenchmarkVerifyProof10000(b *testing.B) {
// 	trie, vals := randomTrie(10000)
// 	root := trie.Hash()
// 	var keys []string
// 	var proofs []*ethdb.MemDatabase
// 	for k := range vals {
// 		keys = append(keys, k)
// 		proof := ethdb.NewMemDatabase()
// 		trie.Prove([]byte(k), 0, proof)
// 		proofs = append(proofs, proof)
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		im := i % len(keys)
// 		if _, err, _ := VerifyProof(root, []byte(keys[im]), proofs[im]); err != nil {
// 			b.Fatalf("key %x: %v", keys[im], err)
// 		}
// 	}
// }
//
// func BenchmarkVerifyProof100000(b *testing.B) {
// 	trie, vals := randomTrie(100000)
// 	root := trie.Hash()
// 	var keys []string
// 	var proofs []*ethdb.MemDatabase
// 	for k := range vals {
// 		keys = append(keys, k)
// 		proof := ethdb.NewMemDatabase()
// 		trie.Prove([]byte(k), 0, proof)
// 		proofs = append(proofs, proof)
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		im := i % len(keys)
// 		if _, err, _ := VerifyProof(root, []byte(keys[im]), proofs[im]); err != nil {
// 			b.Fatalf("key %x: %v", keys[im], err)
// 		}
// 	}
// }

func randomTrie(n int) (*Trie, map[string]*kv) {
	trie := new(Trie)
	vals := make(map[string]*kv)
	for i := byte(0); i < 100; i++ {
		value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		value2 := &kv{common.LeftPadBytes([]byte{i + 10}, 32), []byte{i}, false}
		trie.Update(value.k, value.v)
		trie.Update(value2.k, value2.v)
		vals[string(value.k)] = value
		vals[string(value2.k)] = value2
	}
	for i := 0; i < n; i++ {
		value := &kv{randBytes(32), randBytes(20), false}
		trie.Update(value.k, value.v)
		vals[string(value.k)] = value
	}
	return trie, vals
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	crand.Read(r)
	return r
}
