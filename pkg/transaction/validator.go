package transaction

import (
	"fmt"
	"math/big"
)

// Validate verifies that the Context is in a valid state.
func (ctx *Context) Validate() error {
	summary := ctx.GetAssetMap()

	for asset, quantity := range summary {
		if quantity.Cmp(big.NewRat(0, 1)) == -1 {
			q, _ := quantity.Float32()
			return fmt.Errorf("negative asset quantity %s: %f < 0", asset.Id, q)
		}
	}

	return nil
}
