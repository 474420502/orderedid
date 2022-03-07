package orderedid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"
)

var msgStartTimeUnix uint64 = func() uint64 {
	t, err := time.Parse("2006-01-02", "2022-03-04")
	if err != nil {
		panic(err)
	}
	return uint64(t.UnixMilli())
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
}

const timestampBits uint64 = 64 - 45
const nodeidBits uint64 = 10
const nodeidMark uint64 = 0b1111111111 //  1bit of 10

// New nodeid <= 1024
func New(nodeid uint16) *OrderedIDCreator {

	if nodeid >= (1 << nodeidBits) {
		panic(fmt.Sprintf("nodeid must < %d", 1<<nodeidBits))
	}

	creator := &OrderedIDCreator{
		nodeid: uint64(nodeid),
		count:  0,
	}

	return creator
}

// Create Create a OrderID
func (creator *OrderedIDCreator) Create() OrderedID {

	var tid uint64 = uint64(time.Now().UnixMilli())
	tid -= msgStartTimeUnix
	tid = tid << timestampBits
	tid |= (atomic.AddUint64(&creator.count, 1) << nodeidBits)
	tid |= creator.nodeid
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
func (orderedid OrderedID) NodeID() uint16 {
	return uint16(uint64(orderedid) & nodeidMark)
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
