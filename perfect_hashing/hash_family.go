package misc

import (
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
)

// A *universal* hash family
type HashFamily struct {
	Pick func(desiredMax uint64) HashFunction
}

type HashFunction struct {
	Bins  uint64
	Value func(interface{}) uint64
}

func firstPrimeAfter(n uint64) uint64 {
	var i big.Int
	i.SetUint64(n)
	for !i.ProbablyPrime(0) {
		n++
		i.SetUint64(n)
	}
	return n
}

func uintHashFamily(b uint64) HashFamily {
	return HashFamily{func(n uint64) HashFunction {
		p := firstPrimeAfter(n)
		r := int(float64(b) / math.Log2(float64(p)))
		a := make([]uint64, r)
		for i := 0; i < r; i++ {
			a[i] = rand.Uint64()
		}
		h := func(key interface{}) uint64 {
			var x uint64
			h := uint64(0)
			switch v := key.(type) {
			case int:
				x = uint64(v)
			case int8:
				x = uint64(v)
			case int16:
				x = uint64(v)
			case int32:
				x = uint64(v)
			case int64:
				x = uint64(v)
			case uint:
				x = uint64(v)
			case uint8:
				x = uint64(v)
			case uint16:
				x = uint64(v)
			case uint32:
				x = uint64(v)
			case uint64:
				x = v
			default:
				panic("bad type: " + reflect.TypeOf(key).String())
			}
			for i := 0; x != 0 && i < r; i++ {
				h += (a[i] * (x % p)) % p
				h %= p
				x /= p
			}
			return h
		}
		return HashFunction{Bins: p, Value: h}
	}}
}

var UintHashFamily HashFamily = uintHashFamily(strconv.IntSize)

var Uint8HashFamily HashFamily = uintHashFamily(8)

var Uint16HashFamily HashFamily = uintHashFamily(16)

var Uint32HashFamily HashFamily = uintHashFamily(32)

var Uint64HashFamily HashFamily = uintHashFamily(64)

var IntHashFamily HashFamily = uintHashFamily(strconv.IntSize)

var Int8HashFamily HashFamily = uintHashFamily(8)

var Int16HashFamily HashFamily = uintHashFamily(16)

var Int32HashFamily HashFamily = uintHashFamily(32)

var Int64HashFamily HashFamily = uintHashFamily(64)
