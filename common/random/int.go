package random

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/exp/constraints"
)

func Int[T constraints.Integer](max T) T {
	number, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return T(number.Int64())
}
