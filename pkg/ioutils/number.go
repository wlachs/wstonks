package ioutils

import (
	"fmt"
	"math/big"
)

// ParseRat converts the number string to big.Rat
func ParseRat(s string) (*big.Rat, error) {
	rat := big.NewRat(0, 1)
	if _, ok := rat.SetString(s); !ok {
		return nil, fmt.Errorf("failed to parse numeric string %s", s)
	}

	return rat, nil
}
