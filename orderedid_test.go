package orderedid

import (
	"log"
	"strconv"
	"testing"
)

func BenchmarkCase(b *testing.B) {
	var id OrderedID
	creator := New(1)
	for i := 0; i < b.N; i++ {
		id = creator.Create()
	}

	b.Log(id)
}

func TestCase1(t *testing.T) {
	for i := 0; i < 1000; i++ {
		var id OrderedID
		creator := New(1)
		id = creator.Create()
		idcmp32, err := ParseBase32(id.Base32())
		if err != nil {
			t.Error(err)
		}

		if id != idcmp32 {
			t.Error("")
		}

		if strconv.FormatUint(uint64(id), 10) != idcmp32.String() {
			t.Error("")
		}

		idcmp64, err := ParseBase58(id.Base58())
		if err != nil {
			t.Error(err)
		}

		if id != idcmp64 {
			t.Error("")
		}

		if strconv.FormatUint(uint64(id), 10) != idcmp64.String() {
			t.Error("")
		}

		if id.NodeID() != 1 {
			t.Error("")
		}
	}
}

func TestCase2(t *testing.T) {

	log.Printf("%b", countMark)

	creator := New(1)
	var m map[uint64]OrderedID = make(map[uint64]OrderedID)

	for i := 0; i < 1000000; i++ {
		var id = creator.Create()

		uid := uint64(id)
		// log.Printf("%b", uid)
		// log.Println(id.Count())

		if oid, ok := m[uid]; ok {
			log.Panicf("len:%d, %0.64b %d %d %d %d", len(m), uid, oid.Timestamp(), id.Timestamp(), id.Count(), oid.Count())
		}
		m[uid] = id
	}
}
