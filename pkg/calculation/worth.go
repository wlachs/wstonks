package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/transaction"
	"math/big"
)

// WorthContext holding data required for calculating asset worth.
type WorthContext struct {
	AssetContext       *asset.Context
	TransactionContext *transaction.Context
}

// GetAssetWorthMap calculates the current worth of all assets contained in the transaction.Context with the help of the live asset values
// retrieved from the asset.Context and maps the assets to their corresponding current worth.
func GetAssetWorthMap(ctx WorthContext) (map[*asset.Asset]big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return nil, fmt.Errorf("asset or transaction context is missing")
	}

	assetTxMap := transaction.GetAssetKeyMap(*txCtx)
	assetPriceMap := asset.GetAssetKeyPriceMap(*assetCtx)
	assetMap := asset.GetAssetKeyMap(*assetCtx)
	m := map[*asset.Asset]big.Rat{}

	for a, quantity := range assetTxMap {
		up, ok := assetPriceMap[a]
		if !ok {
			return nil, fmt.Errorf("no unit price found for asset %s", a)
		}

		unitPrice := big.NewRat(0, 1).Set(&up)
		m[assetMap[a]] = *unitPrice.Mul(unitPrice, &quantity)
	}

	return m, nil
}

// GetAssetWorth calculates the current worth of all assets contained in the transaction.Context with the help of the live asset values
// retrieved from the asset.Context.
func GetAssetWorth(ctx WorthContext) (big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return big.Rat{}, fmt.Errorf("asset or transaction context is missing")
	}

	m, err := GetAssetWorthMap(ctx)
	if err != nil {
		return big.Rat{}, err
	}

	worth := big.NewRat(0, 1)
	for _, w := range m {
		worth.Add(worth, &w)
	}

	return *worth, nil
}
