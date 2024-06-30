package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"math/big"
)

// GetAssetReturnMap calculates the current return on all quantities contained in the transaction.Context with the help of the live asset
// values retrieved from the asset.Context and maps the quantities to their corresponding current return.
func (ctx *Context) GetAssetReturnMap() (map[*asset.Asset]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	txCtx := ctx.TransactionContext

	if assetCtx == nil || txCtx == nil {
		return nil, fmt.Errorf("asset or transaction context is missing")
	}

	assetWorthMap, err := ctx.GetAssetWorthMap()
	if err != nil {
		return nil, err
	}

	m := map[*asset.Asset]*big.Rat{}
	assetInitialMap := txCtx.GetAssetKeyInitialWorthMap()

	for a, currentWorth := range assetWorthMap {
		initialWorth, ok := assetInitialMap[a.Id]
		if ok {
			ret := big.NewRat(0, 1)
			ret.Set(currentWorth)
			ret.Sub(ret, initialWorth)
			m[a] = ret
		}
	}

	return m, nil
}
