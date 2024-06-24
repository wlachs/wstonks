package transaction

import (
	"math/big"
)

// GetAssetMap calculates the overall owned assets based on the Context.
// Assets that can be found in the transaction history but have already been sold will not be shown in the summary.
func GetAssetMap(ctx Context) map[*Asset]big.Rat {
	summary := map[*Asset]big.Rat{}

	for i := range ctx.Assets {
		asset := ctx.Assets[i]
		quantity := big.NewRat(0, 1)

		for _, transaction := range asset.Transactions {
			switch transaction.Type {
			case BUY:
				quantity.Add(quantity, &transaction.Quantity)
			case SELL:
				quantity.Sub(quantity, &transaction.Quantity)
			default:
				// not relevant
			}
		}

		if quantity.Cmp(big.NewRat(0, 1)) != 0 {
			summary[asset] = *quantity
		}
	}

	return summary
}

// GetAssetKeyMap calculates the overall owned assets based on the Context while using only the Asset ID as key.
func GetAssetKeyMap(ctx Context) map[string]big.Rat {
	keyMap := map[string]big.Rat{}
	summary := GetAssetMap(ctx)

	for asset, quantity := range summary {
		keyMap[asset.Id] = quantity
	}

	return keyMap
}
