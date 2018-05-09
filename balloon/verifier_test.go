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

package balloon

import (
	"testing"

	"github.com/bbva/qed/balloon/hashing"
	"github.com/bbva/qed/balloon/history"
	"github.com/bbva/qed/balloon/hyper"
	"github.com/bbva/qed/balloon/storage/cache"
)

type FakeVerifiable struct {
	result bool
}

func NewFakeVerifiable(result bool) *FakeVerifiable {
	return &FakeVerifiable{result}
}

func (f *FakeVerifiable) Verify(commitment, event []byte, version uint64) bool {
	return f.result
}

func TestVerify(t *testing.T) {

	testCases := []struct {
		exists         bool
		hyperOK        bool
		historyOK      bool
		queryVersion   uint64
		actualVersion  uint64
		expectedResult bool
	}{
		// Event exists, queryVersion <= actualVersion, and both trees verify it
		{true, true, true, uint64(0), uint64(0), true},
		// Event exists, queryVersion <= actualVersion, but HyperTree does not verify it
		{true, false, true, uint64(0), uint64(0), false},
		// Event exists, queryVersion <= actualVersion, but HistoryTree does not verify it
		{true, true, false, uint64(0), uint64(0), false},

		// Event exists, queryVersion > actualVersion, and both trees verify it
		{true, true, true, uint64(1), uint64(0), true},
		// Event exists, queryVersion > actualVersion, but HyperTree does not verify it
		{true, false, true, uint64(1), uint64(0), false},

		// Event does not exist, HyperTree verifies it
		{false, true, false, uint64(0), uint64(0), true},
		// Event does not exist, HyperTree does not verify it
		{false, false, false, uint64(0), uint64(0), false},
	}

	for i, c := range testCases {
		event := []byte("Yadda yadda")
		commitment := &Commitment{
			[]byte("Some hyperDigest"),
			[]byte("Some historyDigest"),
			c.actualVersion,
		}
		proof := NewProof(
			c.exists,
			NewFakeVerifiable(c.hyperOK),
			NewFakeVerifiable(c.historyOK),
			c.queryVersion,
			c.actualVersion,
			hashing.XorHasher,
		)
		result := proof.Verify(commitment, event)

		if result != c.expectedResult {
			t.Fatalf("Unexpected result '%v' in test case '%d'", result, i)
		}
	}
}

func createBalloon(id string, hasher hashing.Hasher) (*HyperBalloon, func()) {
	frozen, frozenCloseF := openBPlusStorage()
	leaves, leavesCloseF := openBPlusStorage()
	cache := cache.NewSimpleCache(0)

	hyperT := hyper.NewFakeTree(string(0x0), cache, leaves, hasher)
	historyT := history.NewFakeTree(frozen, hasher)
	balloon := NewHyperBalloon(hasher, historyT, hyperT)

	return balloon, func() {
		frozenCloseF()
		leavesCloseF()
	}
}

func TestAddAndVerify(t *testing.T) {
	id := string(0x0)
	hasher := hashing.Sha256Hasher

	balloon, closeF := createBalloon(id, hasher)
	defer closeF()

	key := []byte("Never knows best")
	// keyDigest := hasher(key)

	commitment := <-balloon.Add(key)
	membershipProof := <-balloon.GenMembershipProof(key, commitment.Version)

	historyProof := history.NewFakeProof(membershipProof.HistoryProof, commitment.Version, hasher)
	hyperProof := hyper.NewFakeProof(id, membershipProof.HyperProof, hasher)

	proof := NewProof(
		membershipProof.Exists,
		hyperProof,
		historyProof,
		membershipProof.QueryVersion,
		membershipProof.ActualVersion,
		hasher,
	)

	correct := proof.Verify(commitment, key)

	if !correct {
		t.Errorf("Proof is incorrect")
	}

}
