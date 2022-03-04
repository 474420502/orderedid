package orderedid

import (
	"testing"
	"time"
)

func BenchmarkCaseMID(b *testing.B) {

	var id OrderedID
	creator := New(1)
	for i := 0; i < b.N; i++ {
		id = creator.Create()
	}

	b58 := id.Base58()
	b32 := id.Base32()
	b.Log(time.UnixMilli(int64(id.Timestamp())), id.Base58(), id.Base32(), uint64(id), id)
	b.Log(ParseBase58(b58))
	b.Log(ParseBase32(b32))
}

func TestCase1(t *testing.T) {

	for i := 0; i < 1000; i++ {
		var id OrderedID
		creator := New(1)
		id = creator.Create()
		idcmp32, err := ParseBase32(id.Base32())
		if err != nil {
			panic(err)
		}

		if id != idcmp32 {
			panic("")
		}

		idcmp64, err := ParseBase32(id.Base32())
		if err != nil {
			panic(err)
		}

		if id != idcmp64 {
			panic("")
		}
	}
}
