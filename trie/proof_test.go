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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

func init() {
	mrand.Seed(time.Now().Unix())
}

func TestProof(t *testing.T) {
	trie, vals := randomTrie(500)
	root := trie.Hash()
	for _, kv := range vals {
		proofs := ethdb.NewMemDatabase()
		if trie.Prove(kv.k, 0, proofs) != nil {
			t.Fatalf("missing key %x while constructing proof", kv.k)
		}
		val, err, _ := VerifyProof(root, kv.k, proofs)
		if err != nil {
			t.Fatalf("VerifyProof error for key %x: %v\nraw proof: %v", kv.k, err, proofs)
		}
		if !bytes.Equal(val, kv.v) {
			t.Fatalf("VerifyProof returned wrong value for key %x: got %x, want %x", kv.k, val, kv.v)
		}
	}
}

func TestProof2(t *testing.T) {
	trie, vals := randomTrie(500)
	root := trie.Hash()
	for _, kv := range vals {
		proofs, err := trie.Prove2(kv.k, 0)
		if err != nil {
			t.Fatalf("missing key %x while constructing proof", kv.k)
		}

		// print out proof
		fmt.Println("hi")
		fmt.Println(proofs)
		total := 0
		for _, p := range proofs {
			fmt.Println(p)
			fmt.Println(len(p))
			total += len(p)
		}
		fmt.Println(total)
		fmt.Println("ho")


		val, err, _ := VerifyProof2(root, kv.k, proofs)
		if err != nil {
			t.Fatalf("VerifyProof error for key %x: %v\nraw proof: %v", kv.k, err, proofs)
		}
		if !bytes.Equal(val, kv.v) {
			t.Fatalf("VerifyProof returned wrong value for key %x: got %x, want %x", kv.k, val, kv.v)
		} else {
			fmt.Printf("val: %+v\n", val)
			fmt.Printf("kv.v: %+v\n", kv.v)
		}

		t.Fatalf("uhhhh")
	}
}

func TestOneElementProof(t *testing.T) {
	trie := new(Trie)
	updateString(trie, "k", "v")
	proofs := ethdb.NewMemDatabase()
	trie.Prove([]byte("k"), 0, proofs)
	if len(proofs.Keys()) != 1 {
		t.Error("proof should have one element")
	}
	val, err, _ := VerifyProof(trie.Hash(), []byte("k"), proofs)
	if err != nil {
		t.Fatalf("VerifyProof error: %v\nproof hashes: %v", err, proofs.Keys())
	}
	if !bytes.Equal(val, []byte("v")) {
		t.Fatalf("VerifyProof returned wrong value: got %x, want 'k'", val)
	}
}

func TestVerifyBadProof(t *testing.T) {
	trie, vals := randomTrie(800)
	root := trie.Hash()
	for _, kv := range vals {
		proofs := ethdb.NewMemDatabase()
		trie.Prove(kv.k, 0, proofs)
		if len(proofs.Keys()) == 0 {
			t.Fatal("zero length proof")
		}
		keys := proofs.Keys()
		key := keys[mrand.Intn(len(keys))]
		node, _ := proofs.Get(key)
		proofs.Delete(key)
		mutateByte(node)
		proofs.Put(crypto.Keccak256(node), node)
		if _, err, _ := VerifyProof(root, kv.k, proofs); err == nil {
			t.Fatalf("expected proof to fail for key %x", kv.k)
		}
	}
}

// // Why did I use an iterator instead of looping through vals??
// func TestGetProofNew(t *testing.T) {
// 	trie, vals := randomTrie(500)
// 	for _, kv := range vals {
// 		fmt.Printf("kv.k: %+v\n", kv.k)
// 		fmt.Printf("kv.v: %+v\n", kv.v)
// 		proof, pkey := trie.GetProof(kv.k)
// 		fmt.Printf("proof: %+v\n", proof)
// 		fmt.Printf("pkey: %+v\n", pkey)
//
// 		// need to get node at that key in order to check
// 		// HASH KV.V (is that the RLP????)
// 		// GET RLP FROM KV...
//
// 		// match, proot := VerifyProofOf(n, trie.root, proof, pkey)
// 		// if !match {
// 		// 	t.Fatalf("VerifyProofOf returned wrong value for key %x: got %x, want %x", it.path, proot, trie.root)
// 		// } else {
// 		// 	t.Errorf("VerifyProofOf returned the RIGHT value for key %x: got %x, want %x", it.path, proot, trie.root)
// 		// }
// 	}
// }


// func TestGetProofIterator(t *testing.T) {
// 	db, trie, _ := makeTestTrie()
// 	// Gather all the node hashes found by the iterator
// 	for it := trie.NodeIterator(nil); it.Next(true); {
// 		if it.Hash() != (common.Hash{}) {
// 			fmt.Printf("HASH: %+v\n", it.Hash())
// 			fmt.Printf("PATH: %+v\n", it.Path())
//
// 			val, _ := trie.TryGet(it.Path())
// 			fmt.Printf("val: %+v\n", val)
//
// 			proof, pkey := trie.GetProof(it.Path())
// 			fmt.Printf("proof: %x\n", proof)
// 			fmt.Printf("pkey: %+v\n", pkey)
//
// 			n, preimage := db.Node(it.Hash())
// 			fmt.Printf("node: %+v\n", n)
// 			fmt.Printf("preimage: %+v\n", preimage)
//
// 			// need to get actual node...from db maybe
// 			match, proot := VerifyProofOf(n, trie.root, proof, pkey)
// 		}
// 	}
// }

// TestGetProof constructs a proof
// func TestGetProof(t *testing.T) {
// 	trie, _ := randomTrie(500)
// 	// it := NewIterator(trie.NodeIterator(nil))
// 	// use iterator to iterate through nodes
// 	it := &nodeIterator{trie: trie} // create iterator for trie
// 	for it.Next(true) {
// 		// get node iterator is on
// 		itstate, _, _, _ := it.peek(false) // FOR SOME REASON THIS IS <NIL>
// 		fmt.Printf("itstate: %+v\n", itstate)
// 		n := itstate.node
// 		if n != nil {
// 			proof, pkey := trie.GetProof(it.path)
// 			fmt.Printf("it.Key: %+v\n", it.path)
// 			fmt.Printf("node: %x\n", n)
// 			match, proot := VerifyProofOf(n, trie.root, proof, pkey)
// 			if !match {
// 				t.Fatalf("VerifyProofOf returned wrong value for key %x: got %x, want %x", it.path, proot, trie.root)
// 			} else {
// 				t.Errorf("VerifyProofOf returned the RIGHT value for key %x: got %x, want %x", it.path, proot, trie.root)
// 			}
// 		}
// 	}
// }

// TEST THAT WE'RE CALCULATING PARENT HASHES CORRECTLY
// func TestParentHashing(t *testing.T) {
// 	trie := newEmpty()
//
// 	trie.Update([]byte("hi"), []byte("there"))
// 	trie.Update([]byte("howyou"), []byte("doin"))
// 	trie.Update([]byte("mmm"), []byte("hmmm"))
// 	// parentHash := trie.Hash()
//
// 	trie.Decode(trie.root, 4)
// 	t.Fatalf("wtfffffff\n")
//
// 	// enc, _ = rlp.EncodeToBytes(trie.root)
// }

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
			proofs := ethdb.NewMemDatabase()
			if trie.Prove(kv.k, 0, proofs); len(proofs.Keys()) == 0 {
				b.Fatalf("zero length proof for %x", kv.k)
			}

			averageProofsize += (len(proofs.Keys()))
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
