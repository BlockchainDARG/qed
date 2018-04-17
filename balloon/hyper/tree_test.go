// Copyright © 2018 Banco Bilbao Vizcaya Argentaria S.A.  All rights reserved.
// Use of this source code is governed by an Apache 2 License
// that can be found in the LICENSE file

package hyper

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
	"verifiabledata/balloon/hashing"
	"verifiabledata/balloon/storage/badger"
	"verifiabledata/balloon/storage/bplus"
	"verifiabledata/balloon/storage/cache"
)

func TestAdd(t *testing.T) {
	store, closeF := openBPlusStorage()
	defer closeF()

	cache := cache.NewSimpleCache(5000)
	hasher := hashing.XorHasher

	ht := NewTree(string(0x0), 2, cache, store, hasher, fakeLeafHasherF(hasher), fakeInteriorHasherF(hasher))

	key := []byte{0x5a}
	value := []byte{0x01}

	expectedCommitment := []byte{0x5a}
	commitment := <-ht.Add(key, value)

	if bytes.Compare(commitment, expectedCommitment) != 0 {
		t.Fatalf("Expected: %x, Actual: %x", expectedCommitment, commitment)
	}

}

func TestAuditPath(t *testing.T) {

	store, closeF := openBPlusStorage()
	defer closeF()

	cache := cache.NewSimpleCache(5000)
	hasher := hashing.XorHasher

	ht := NewTree(string(0x0), 2, cache, store, hasher, fakeLeafHasherF(hasher), fakeInteriorHasherF(hasher))

	key := []byte{0x5a}
	value := []byte{0x01}

	<-ht.Add(key, value)

	expectedPath := [][]byte{
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
		[]byte{0x00},
	}
	proof := <-ht.AuditPath(key)

	if !comparePaths(expectedPath, proof.AuditPath) {
		t.Fatalf("Invalid path: expected %v, actual %v", expectedPath, proof.AuditPath)
	}

}

func comparePaths(expected, actual [][]byte) bool {
	for i, e := range expected {
		if !bytes.Equal(e, actual[i]) {
			return false
		}
	}
	return true
}

func randomBytes(n int) []byte {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes
}

func BenchmarkAdd(b *testing.B) {
	store, closeF := openBadgerStorage()
	defer closeF()
	cache := cache.NewSimpleCache(5000000)
	hasher := hashing.Sha256Hasher
	ht := NewTree("my test tree", 30, cache, store, hasher, LeafHasherF(hasher), InteriorHasherF(hasher))
	b.N = 100000
	for i := 0; i < b.N; i++ {
		key := randomBytes(64)
		value := randomBytes(1)
		store.Add(key, value)
		<-ht.Add(key, value)
	}
	b.Logf("stats = %+v\n", ht.stats)
}

func openBPlusStorage() (*bplus.BPlusTreeStorage, func()) {
	store := bplus.NewBPlusTreeStorage()
	return store, func() {
		store.Close()
	}
}

func openBadgerStorage() (*badger.BadgerStorage, func()) {
	store := badger.NewBadgerStorage("/tmp/hyper_tree_test.db")
	return store, func() {
		store.Close()
		deleteFile("/tmp/hyper_tree_test.db")
	}
}

func deleteFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Unable to remove db file %s", err)
	}
}
