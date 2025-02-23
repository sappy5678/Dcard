package shortcodeid

import (
	"errors"
	"math"
	"sync/atomic"
)

// base57, ignore 0, 1, I, O, l
const base57Chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var base57Map map[byte]uint64

func init() {
	base57Map = make(map[byte]uint64)
	for i, c := range base57Chars {
		base57Map[byte(c)] = uint64(i)
	}
}

type impl struct {
	machineID uint64
	counter   uint64
}

// New returns a new Repository
func New(machineID uint64) Repository {
	return &impl{
		machineID: machineID,
		counter:   0,
	}
}

// NextID returns a new short code
func (im *impl) NextID() string {
	machineID := im.encodeID(im.machineID)
	newCounter := atomic.AddUint64(&im.counter, 1)
	counterID := im.encodeID(newCounter)
	return machineID + "-" + counterID
}

func (im *impl) encodeID(n uint64) string {
	if n == 0 {
		return string(base57Chars[0])
	}

	length := uint64(len(base57Chars))
	var encoded []byte
	for n > 0 {
		r := n % length
		n /= length
		encoded = append([]byte{base57Chars[r]}, encoded...)
	}
	return string(encoded)
}

// decodeID decodes base57 string to number
// in there, just for testing
func (im *impl) decodeID(encoded string) (uint64, error) {
	var number uint64
	length := uint64(len(base57Chars))

	for i, c := range []byte(encoded) {
		val, exists := base57Map[c]
		if !exists {
			return 0, errors.New("invalid character in encoded string")
		}
		number += val * uint64(math.Pow(float64(length), float64(len(encoded)-i-1)))
	}

	return number, nil
}
