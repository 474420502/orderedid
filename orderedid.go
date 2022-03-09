package orderedid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var msgStartTimeUnix uint64 = func() uint64 {
	t, err := time.Parse("2006-01-02", "2022-03-04")
	if err != nil {
		panic(err)
	}
	return uint64(t.UnixNano() / 1000000)
}()

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// ErrInvalidBase58 is returned by ParseBase58 when given an invalid []byte
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 is returned by ParseBase32 when given an invalid []byte
var ErrInvalidBase32 = errors.New("invalid base32")

// Create maps for decoding Base58/Base32.
// This speeds up the process tremendously.
func init() {

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// OrderedID 有序的id
type OrderedID uint64

type OrderedIDCreator struct {
	nodeid uint64
	count  uint64

	lock sync.Mutex
}

const timestampBits uint64 = 64 - 43            // 43bit 支持 2022-03-04 后  200多年时间
const nodeidBits uint64 = 5                     // 21 - 5 = 16 bit
const nodeidMark uint64 = (1 << nodeidBits) - 1 //  1bit of 5
const countlimit uint64 = 1 << (timestampBits - nodeidBits)
const countMark uint64 = (countlimit - 1) << 5

// New nodeid < 32
func New(nodeid uint8) *OrderedIDCreator {

	if uint64(nodeid) > nodeidMark {
		panic(fmt.Sprintf("nodeid must < %d", 1<<nodeidBits))
	}

	creator := &OrderedIDCreator{
		nodeid: uint64(nodeid),
		count:  countlimit,
	}

	return creator
}

// Create Create a OrderID
func (creator *OrderedIDCreator) Create() OrderedID {

	var tid uint64 = uint64(time.Now().UnixNano()) / 100000
	creator.lock.Lock()
	if creator.count >= countlimit {
		creator.count = 0
	}
	count := creator.count
	creator.count++
	creator.lock.Unlock()

	tid -= msgStartTimeUnix      // 减去相对时间
	tid = (tid << timestampBits) // 偏移到占用位
	tid |= (count << nodeidBits) //
	tid |= creator.nodeid        //

	return OrderedID(tid)

}

// Bytes return return the  bytes(the 8 byte of int64) of ordererid
func (orderedid OrderedID) Bytes() []byte {
	var buf []byte = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(orderedid))
	return buf
}

// String return the number string
func (orderedid OrderedID) String() string {
	return strconv.FormatUint(uint64(orderedid), 10)
}

// Timestamp return the timestamp
func (orderedid OrderedID) Timestamp() uint64 {
	return (uint64(orderedid) >> timestampBits) + msgStartTimeUnix
}

// NodeID return the NodeID
func (orderedid OrderedID) NodeID() uint64 {
	return uint64(orderedid) & nodeidMark
}

// NodeID return the NodeID
func (orderedid OrderedID) Count() uint64 {
	return (uint64(orderedid) & countMark) >> 5
}

// Base32 return a base32 string
func (orderedid OrderedID) Base32() string {

	if orderedid < 32 {
		return string(encodeBase32Map[orderedid])
	}

	b := make([]byte, 0, 12)
	for orderedid >= 32 {
		b = append(b, encodeBase32Map[orderedid%32])
		orderedid /= 32
	}
	b = append(b, encodeBase32Map[orderedid])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// Uint64 return the  int64 , orderedid
func (orderedid OrderedID) Uint64() uint64 {
	return uint64(orderedid)
}

// Base58 return a base58 string
func (orderedid OrderedID) Base58() string {

	if orderedid < 58 {
		return string(encodeBase58Map[orderedid])
	}

	b := make([]byte, 0, 11)
	for orderedid >= 58 {
		b = append(b, encodeBase58Map[orderedid%58])
		orderedid /= 58
	}
	b = append(b, encodeBase58Map[orderedid])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseUint64 converts an uint64 into a OrderedID
func ParseUint64(id uint64) OrderedID {
	return OrderedID(id)
}

// ParseString converts a string into a OrderedID
func ParseString(id string) (OrderedID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return OrderedID(i), err
}

// ParseBase32 parses a base32 []byte into a OrderedID
func ParseBase32(b32 string) (OrderedID, error) {

	var id uint64

	for _, char := range *(*[]byte)(unsafe.Pointer(&b32)) {
		if decodeBase32Map[char] == 0xFF {
			return 0, ErrInvalidBase32
		}
		id = id*32 + uint64(decodeBase32Map[char])
	}

	return OrderedID(id), nil
}

// ParseBase58 parses a base58 []byte into a snowflake ID
func ParseBase58(b58 string) (OrderedID, error) {
	var id uint64
	for _, char := range *(*[]byte)(unsafe.Pointer(&b58)) {
		if decodeBase58Map[char] == 0xFF {
			return 0, ErrInvalidBase58
		}
		id = id*58 + uint64(decodeBase58Map[char])
	}
	return OrderedID(id), nil
}
