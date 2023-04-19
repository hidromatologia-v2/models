package random

import (
	cryptoRand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Float[T constraints.Float](max T) T {
	bi, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
	return T(rand.New(rand.NewSource(bi.Int64())).Float64()) * max
}
