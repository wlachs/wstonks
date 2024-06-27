package asset

import "math/big"

// GetAssetKeyPriceMap creates a map mapping quantities to their respective live unit prices.
func (ctx *Context) GetAssetKeyPriceMap() map[string]*big.Rat {
	m := map[string]*big.Rat{}

	for _, asset := range ctx.Assets {
		m[asset.Id] = asset.UnitPrice
	}

	return m
}
