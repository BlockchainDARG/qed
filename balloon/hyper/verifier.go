// Copyright © 2018 Banco Bilbao Vizcaya Argentaria S.A.  All rights reserved.
// Use of this source code is governed by an Apache 2 License
// that can be found in the LICENSE file

package hyper

import (
	"bytes"
	"log"
	"os"
	"verifiabledata/balloon/hashing"
)

type Verifier struct {
	id             []byte
	digestLength   int
	leafHasher     LeafHasher
	interiorHasher InteriorHasher
	log            *log.Logger
}

func NewVerifier(id string, hasher hashing.Hasher, leafHasher LeafHasher, interiorHasher InteriorHasher) *Verifier {
	digestLength := len(hasher([]byte("x"))) * 8
	return &Verifier{
		[]byte(id),
		digestLength,
		leafHasher,
		interiorHasher,
		log.New(os.Stdout, "HistoryVerifier", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}
}

func (v *Verifier) Verify(expectedDigest []byte, auditPath [][]byte, key, value []byte) (bool, []byte) {
	recomputed := v.rootHash(auditPath, rootPosition(v.digestLength), key, value)
	return bytes.Equal(expectedDigest, recomputed), recomputed
}

func (v *Verifier) rootHash(auditPath [][]byte, pos *Position, key, value []byte) []byte {
	if pos.height == 0 {
		return v.leafHasher(v.id, value, pos.base)
	}
	if !bitIsSet(key, v.digestLength-pos.height) { // if k_j == 0
		left := v.rootHash(auditPath, pos.left(), key, value)
		right := auditPath[pos.height]
		next := pos.right()
		return v.interiorHasher(left, right, next.base, next.heightBytes())
	}
	left := auditPath[pos.height]
	right := v.rootHash(auditPath, pos.right(), key, value)
	next := pos.left()
	return v.interiorHasher(left, right, next.base, next.heightBytes())
}
