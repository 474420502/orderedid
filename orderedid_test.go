package orderedid

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCase(b *testing.B) {
	var id OrderedID
	creator := New()
	for i := 0; i < b.N; i++ {
		id = creator.Create()
	}

	b.Log(id)
}

func BenchmarkCase2(b *testing.B) {

}

func TestCase1(t *testing.T) {
	creator := New()
	for i := 0; i < 1000; i++ {

		var id OrderedID = creator.Create()

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

		idcmp58, err := ParseBase58(id.Base58())
		if err != nil {
			t.Error(err)
		}

		if id != idcmp58 {
			t.Error("")
		}

		if strconv.FormatUint(uint64(id), 10) != idcmp58.String() {
			t.Error("")
		}

		idcmp64, err := ParseBase64(id.Base64())
		if err != nil {
			t.Error(err)
		}

		if id != idcmp64 {
			t.Error("")
		}

		if strconv.FormatUint(uint64(id), 10) != idcmp64.String() {
			t.Error("")
		}
	}
}

func TestCaseBitsFullPanic(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("should be Error")
		}
	}()

	var id OrderedID
	for i := 0; i < 2; i++ {
		creator := NewWith(1)
		id = creator.Create()
		// log.Println(id.Base64(), id.String())
	}

	t.Error(id)
}

func TestCase2(t *testing.T) {

	log.Printf("%b", countMark)

	creator := NewWith(2)
	defer creator.Destroy()

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
		if id.NodeID() != 2 {
			t.Error("")
		}
	}
}

func TestCaseBitsFullPanic2(t *testing.T) {
	time.Sleep(time.Second * 1)

	defer func() {
		if err := recover(); err == nil {
			t.Error("should be Error")
		}
	}()

	var id OrderedID
	for i := 0; i < 100; i++ {
		creator := New()
		id = creator.Create()

		// log.Println(id.Base64(), id.String())
	}

	t.Error(id)
}

func TestBase32(t *testing.T) {
	// ORUGS4TFMNQXEBLVBA2W35P7CI6

	creator := New()
	log.Println(creator.Create().Base58())
}
