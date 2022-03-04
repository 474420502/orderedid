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

func BenchmarkCaseMID1(b *testing.B) {

}

func TestCase12312(t *testing.T) {
	t.Errorf("%0.64b", ^uint64(18446462598732840960))

}
