// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package badger

import (
	"bytes"
	"github.com/bbva/qed/balloon/storage"

	b "github.com/dgraph-io/badger"
	"github.com/bbva/qed/log"
)

type BadgerStorage struct {
	db *b.DB
}

func (s *BadgerStorage) Add(key []byte, value []byte) error {
	return s.db.Update(func(txn *b.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *BadgerStorage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(txn *b.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}
		return nil
	})
	switch err {
	case nil:
		return value, nil
	case b.ErrEmptyKey:
		return make([]byte, 0), nil
	default:
		return nil, err
	}
}

func (s *BadgerStorage) GetRange(start, end []byte) storage.LeavesSlice {
	var leaves storage.LeavesSlice

	s.db.View(func(txn *b.Txn) error {
		opts := b.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(start); it.Valid(); it.Next() {
			item := it.Item()
			var k []byte
			k = item.KeyCopy(k)
			if bytes.Compare(k, end) > 0 {
				break
			}
			leaves = append(leaves, k)
		}
		return nil
	})

	return leaves
}

func (s *BadgerStorage) Delete(key []byte) error {
	return s.db.Update(func(txn *b.Txn) error {
		return txn.Delete(key)
	})
}

func (s *BadgerStorage) Close() error {
	return s.db.Close()
}

func NewBadgerStorage(path string) *BadgerStorage {
	opts := b.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	opts.SyncWrites = false
	db, err := b.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &BadgerStorage{db}
}
