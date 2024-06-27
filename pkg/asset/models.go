package asset

import "math/big"

// Asset holds detailed information about an asset
type Asset struct {
	Id        string
	UnitPrice *big.Rat
}
