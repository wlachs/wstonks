package ioutils

import (
	"fmt"
	"math/big"
)

// ParseRat converts the number string to big.Rat
func ParseRat(s string) (big.Rat, error) {
	rat := big.Rat{}
	if _, ok := rat.SetString(s); !ok {
		return big.Rat{}, fmt.Errorf("failed to parse numeric string %s", s)
	}

	return rat, nil
}
