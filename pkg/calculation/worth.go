package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"math/big"
)

// GetAssetWorthMap calculates the current worth of all quantities contained in the transaction.Context with the help of the live asset values
// retrieved from the asset.Context and maps the quantities to their corresponding current worth.
func (ctx *Context) GetAssetWorthMap() (map[*asset.Asset]*big.Rat, error) {
	assetCtx := ctx.AssetContext

	if assetCtx == nil {
		return nil, fmt.Errorf("asset context is missing")
	}

	assets := assetCtx.Assets
	return ctx.GetAssetWorthMapOfAssets(assets)
}

// GetAssetWorthMapOfAssets calculates the current worth of all quantities contained in the transaction.Context of the given assets with the
// help of the live asset values retrieved from the asset.Context and maps the quantities to their corresponding current worth.
func (ctx *Context) GetAssetWorthMapOfAssets(assets []*asset.Asset) (map[*asset.Asset]*big.Rat, error) {
	assetKeys := make([]string, 0, len(assets))
	for _, a := range assets {
		assetKeys = append(assetKeys, a.Id)
	}

	assetWorthMap, err := ctx.GetAssetKeyWorthMapOfKeys(assetKeys)
	if err != nil {
		return nil, err
	}

	m := map[*asset.Asset]*big.Rat{}
	for _, a := range assets {
		m[a] = assetWorthMap[a.Id]
	}

	return m, nil
}

// GetAssetKeyWorthMapOfKeys calculates the current worth of all quantities contained in the transaction.Context of the given asset keys
// with the help of the live asset values retrieved from the asset.Context and maps the quantities to their corresponding current worth.
func (ctx *Context) GetAssetKeyWorthMapOfKeys(keys []string) (map[string]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return nil, fmt.Errorf("asset or transaction context is missing")
	}

	assetTxMap := txCtx.GetAssetKeyMap()
	assetPriceMap := assetCtx.GetAssetKeyPriceMap()
	m := map[string]*big.Rat{}

	for _, key := range keys {
		quantity, ok := assetTxMap[key]
		if !ok {
			m[key] = big.NewRat(0, 1)
			continue
		}

		up, ok := assetPriceMap[key]
		if !ok {
			return nil, fmt.Errorf("no unit price found for asset %s", key)
		}

		unitPrice := big.NewRat(0, 1).Set(up)
		m[key] = unitPrice.Mul(unitPrice, quantity)
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
