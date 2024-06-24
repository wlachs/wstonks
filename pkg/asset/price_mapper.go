package asset

import "math/big"

// GetAssetKeyPriceMap creates a map mapping assets to their respective live unit prices.
func GetAssetKeyPriceMap(ctx Context) map[string]big.Rat {
	m := map[string]big.Rat{}

	for _, asset := range ctx.Assets {
		m[asset.Id] = asset.UnitPrice
	}

	return m
}
