// Copyright © 2018 Banco Bilbao Vizcaya Argentaria S.A.  All rights reserved.
// Use of this source code is governed by an Apache 2 License
// that can be found in the LICENSE file

package balloon

import (
	"encoding/binary"
	"verifiabledata/balloon/hashing"
	"verifiabledata/balloon/history"
	"verifiabledata/balloon/hyper"
	"verifiabledata/balloon/storage"
)

type Balloon interface {
	Add(event []byte) chan *Commitment
	Close() chan bool
}

type HyperBalloon struct {
	history *history.Tree
	hyper   *hyper.Tree
	hasher  hashing.Hasher
	version uint
	ops     chan interface{} // serialize operations
}

type Commitment struct {
	HistoryDigest []byte
	IndexDigest   []byte
	Version       uint
}

type AuditPath [][]byte

type MembershipProof struct {
	Exists        bool
	HyperProof    AuditPath
	HistoryProof  AuditPath
	QueryVersion  uint
	ActualVersion uint
}

func NewHyperBalloon(path string, hasher hashing.Hasher, frozen, leaves storage.Store, cache storage.Cache) *HyperBalloon {

	history := history.NewTree(frozen, hasher)
	hyper := hyper.NewTree(path, 30, cache, leaves, hasher, hyper.LeafHasherF(hasher), hyper.InteriorHasherF(hasher))

	b := HyperBalloon{
		history,
		hyper,
		hasher,
		0,
		nil,
	}
	b.ops = b.operations()
	return &b

}

func (b HyperBalloon) Add(event []byte) chan *Commitment {
	result := make(chan *Commitment)
	b.ops <- &add{
		event,
		result,
	}
	return result
}

func (b HyperBalloon) GenMembershipProof(event []byte, version uint) chan *MembershipProof {
	result := make(chan *MembershipProof)
	b.ops <- &membership{
		event,
		version,
		result,
	}
	return result
}

func (b HyperBalloon) Close() chan bool {
	result := make(chan bool)

	b.history.Close()
	b.hyper.Close()

	b.ops <- &close{true, result}
	return result
}

// INTERNALS

type add struct {
	event  []byte
	result chan *Commitment
}

type membership struct {
	event   []byte
	version uint
	result  chan *MembershipProof
}

type close struct {
	stop   bool
	result chan bool
}

// Run listens in channel operations to execute in the tree
func (b *HyperBalloon) operations() chan interface{} {
	operations := make(chan interface{}, 0)
	go func() {
		for {
			select {
			case op := <-operations:
				switch msg := op.(type) {
				case *close:
					msg.result <- true
					return
				case *add:
					digest, _ := b.add(msg.event)
					msg.result <- digest
				case *membership:
					proof, _ := b.genMembershipProof(msg.event, msg.version)
					msg.result <- proof
				default:
					panic("Hyper tree Run() message not implemented!!")
				}

			}
		}
	}()
	return operations
}

func (b *HyperBalloon) add(event []byte) (*Commitment, error) {
	digest := b.hasher(event)
	b.version++
	index := make([]byte, 8)
	binary.LittleEndian.PutUint64(index, uint64(b.version))

	return &Commitment{
		<-b.history.Add(digest, index),
		<-b.hyper.Add(index, digest),
		b.version,
	}, nil
}

func (b *HyperBalloon) genMembershipProof(event []byte, version uint) (*MembershipProof, error) {
	digest := b.hasher(event)

	var hyperProof *hyper.MembershipProof
	var historyProof *[][]byte

	hyperProof = <-b.hyper.Prove(digest)

	var exists bool
	if len(hyperProof.ActualValue) > 0 {
		exists = true
	}

	actualVersion := uint(binary.LittleEndian.Uint64(hyperProof.ActualValue))

	if exists && actualVersion <= version {
		historyProof = <-b.history.Prove(hyperProof.ActualVersion)
	}

	return &MembershipProof{
		exists,
		hyperProof.AuditPath,
		historyProof.AuditPath,
		version,
		actualVersion,
	}, nil

}
