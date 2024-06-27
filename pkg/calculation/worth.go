package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"math/big"
)

// GetAssetWorthMap calculates the current worth of all quantities contained in the transaction.Context with the help of the live asset values
// retrieved from the asset.Context and maps the quantities to their corresponding current worth.
func (ctx *Context) GetAssetWorthMap() (map[*asset.Asset]*big.Rat, error) {
	assetWorthMap, err := ctx.GetAssetKeyWorthMap()
	if err != nil {
		return nil, err
	}

	m := map[*asset.Asset]*big.Rat{}
	assetCtx := ctx.AssetContext
	assetMap := assetCtx.GetAssetKeyMap()

	for a, worth := range assetWorthMap {
		m[assetMap[a]] = worth
	}

	return m, nil
}

// GetAssetKeyWorthMap calculates the current worth of all quantities contained in the transaction.Context with the help of the live asset
// values retrieved from the asset.Context and maps the quantities to their corresponding current worth.
func (ctx *Context) GetAssetKeyWorthMap() (map[string]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return nil, fmt.Errorf("asset or transaction context is missing")
	}

	assetTxMap := txCtx.GetAssetKeyMap()
	assetPriceMap := assetCtx.GetAssetKeyPriceMap()
	m := map[string]*big.Rat{}

	for a, quantity := range assetTxMap {
		up, ok := assetPriceMap[a]
		if !ok {
			return nil, fmt.Errorf("no unit price found for asset %s", a)
		}

		unitPrice := big.NewRat(0, 1).Set(up)
		m[a] = unitPrice.Mul(unitPrice, quantity)
	}

	return m, nil
}

// GetAssetWorth calculates the current worth of all quantities contained in the transaction.Context with the help of the live asset values
// retrieved from the asset.Context.
func (ctx *Context) GetAssetWorth() (*big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return nil, fmt.Errorf("asset or transaction context is missing")
	}

	m, err := ctx.GetAssetWorthMap()
	if err != nil {
		return nil, err
	}

	worth := big.NewRat(0, 1)
	for _, w := range m {
		worth.Add(worth, w)
	}

	return worth, nil
}
