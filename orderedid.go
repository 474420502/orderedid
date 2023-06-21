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

const encodeBase32Map = "xotdrfg8ejkmcpybnq1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

const encodeBase64Map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789$#"

var decodeBase64Map [256]byte

// ErrInvalidBase64 is returned by ParseBase64 when given an invalid []byte
var ErrInvalidBase64 = errors.New("invalid base64")

// ErrInvalidBase58 is returned by ParseBase58 when given an invalid []byte
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 is returned by ParseBase32 when given an invalid []byte
var ErrInvalidBase32 = errors.New("invalid base32")

func init() {

	for i := 0; i < len(encodeBase64Map); i++ {
		decodeBase64Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase64Map); i++ {
		decodeBase64Map[encodeBase64Map[i]] = byte(i)
	}

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

	lastts uint64
	lock   sync.Mutex
}

const timestampBits uint64 = 64 - 43            // 43bit 支持 2022-03-04 后  200多年时间
const nodeidBits uint64 = 6                     // 21 - 6 = 16 bit
const nodeidMark uint64 = (1 << nodeidBits) - 1 //  1bit of 6
const countlimit uint64 = 1 << (timestampBits - nodeidBits)
const countMark uint64 = (countlimit - 1) << nodeidBits

const allusedbits = ^uint64(0)

var newlock sync.Mutex
var nodeidBitsSet uint64 = 0

// NewWith nodeid < 64
func NewWith(nodeid uint8) *OrderedIDCreator {
	if uint64(nodeid) > nodeidMark {
		panic(fmt.Sprintf("nodeid must < %d", 1<<nodeidBits))
	}

	newlock.Lock()
	defer newlock.Unlock()

	checkbit := uint64(1 << nodeid)
	if nodeidBitsSet&checkbit != 0 {
		panic(fmt.Sprintf("nodeid %d is exists.", nodeid))
	}
	nodeidBitsSet |= checkbit // 设置位已经被占用

	creator := &OrderedIDCreator{
		nodeid: uint64(nodeid),
		count:  countlimit,
	}

	return creator
}

// New nodeid < 64
func New() *OrderedIDCreator {

	newlock.Lock()
	defer newlock.Unlock()

	if nodeidBitsSet == allusedbits {
		panic("nodeis is all used")
	}

	var checkbit uint64
	for i := uint64(0); i <= nodeidMark; i++ {
		checkbit = 1 << i
		if nodeidBitsSet&checkbit == 0 {
			nodeidBitsSet |= checkbit // 设置位已经被占用
			creator := &OrderedIDCreator{
				nodeid: i,
				count:  countlimit,
			}

			return creator
		}
	}

	panic("nodeis is full?")
}

// Destroy release the used nodeid.
func (creator *OrderedIDCreator) Destroy() {
	newlock.Lock()
	defer newlock.Unlock()
	nodeidBitsSet &^= (1 << creator.nodeid) // clear bits
}

// Create Create a OrderID
func (creator *OrderedIDCreator) Create() OrderedID {

	for {
		un := uint64(time.Now().UnixNano())
		var tid uint64 = un / 1000000
		creator.lock.Lock()
		if creator.lastts != tid {
			creator.lastts = tid
			creator.count = 0
		} else if creator.count >= countlimit {
			time.Sleep(time.Duration((tid+1)*1000000 - un)) // 如果毫秒位占满, count只能迁移到下个count
			creator.count = 0
			creator.lastts = uint64(time.Now().UnixNano()) / 1000000
			creator.lock.Unlock() // 解锁 循环
			continue
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
	return (uint64(orderedid) & countMark) >> nodeidBits
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

// Base58 return a base64 string
func (orderedid OrderedID) Base64() string {

	if orderedid < 64 {
		return string(encodeBase64Map[orderedid])
	}

	b := make([]byte, 0, 11)
	for orderedid >= 64 {
		b = append(b, encodeBase64Map[orderedid%64])
		orderedid /= 64
	}
	b = append(b, encodeBase64Map[orderedid])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseUint64 converts an uint64 into a OrderedID
func ParseUint64(ordid uint64) OrderedID {
	return OrderedID(ordid)
}

// ParseString converts a string into a OrderedID
func ParseString(ordid string) (OrderedID, error) {
	i, err := strconv.ParseInt(ordid, 10, 64)
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

// ParseBase58 parses a base58 []byte into a OrderedID
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

// ParseBase64 parses a base64 []byte into a OrderedID
func ParseBase64(b64 string) (OrderedID, error) {
	var id uint64
	for _, char := range *(*[]byte)(unsafe.Pointer(&b64)) {
		if decodeBase64Map[char] == 0xFF {
			return 0, ErrInvalidBase64
		}
		id = id*64 + uint64(decodeBase64Map[char])
	}
	return OrderedID(id), nil
}
