package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/transaction"
	"math/big"
	"slices"
)

// GetSalesForReturn calculates how much and which positions should be sold in order to realize the given return.
func (ctx *Context) GetSalesForReturn(r *big.Rat) (map[asset.Asset]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	if assetCtx == nil {
		return nil, fmt.Errorf("asset context not set")
	}

	return ctx.GetSalesForReturnWithAssets(r, ctx.AssetContext.Assets)
}

// GetSalesForReturnWithAssets calculates how much and which positions should be sold in order to realize the given return.
func (ctx *Context) GetSalesForReturnWithAssets(r *big.Rat, assets []asset.Asset) (map[asset.Asset]*big.Rat, error) {
	return nil, nil
}

// GetMaxProfitAndLossForAssets calculates the maximum realizable profit and loss for each asset with live data.
func (ctx *Context) GetMaxProfitAndLossForAssets() (map[*asset.Asset]*big.Rat, map[*asset.Asset]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	if assetCtx == nil {
		return nil, nil, fmt.Errorf("asset context not set")
	}

	profit := map[*asset.Asset]*big.Rat{}
	loss := map[*asset.Asset]*big.Rat{}
	for _, a := range assetCtx.Assets {
		p, l, err := ctx.GetMaxProfitAndLossForAsset(a)
		if err != nil {
			return nil, nil, err
		}
		profit[&a] = p
		loss[&a] = l
	}

	return profit, loss, nil
}

// GetMaxProfitAndLossForAsset checks every open position of the asset and calculates the maximum realizable profit and loss.
func (ctx *Context) GetMaxProfitAndLossForAsset(a asset.Asset) (*big.Rat, *big.Rat, error) {
	txCtx := ctx.TransactionContext
	if txCtx == nil {
		return nil, nil, fmt.Errorf("transaction context not set")
	}

	i := slices.IndexFunc(txCtx.Assets, func(txAsset *transaction.TxAsset) bool {
		return txAsset.Id == a.Id
	})

	if i == -1 {
		return big.NewRat(0, 1), big.NewRat(0, 1), nil
	}

	txAsset := txCtx.Assets[i]
	p := txCtx.GetAssetPositions(txAsset)
	maxProfit, maxLoss := big.NewRat(0, 1), big.NewRat(0, 1)
	diff := big.NewRat(0, 1)

	for _, position := range p {
		diff.Add(diff, position.GetReturnForUnitPrice(a.UnitPrice))
		if maxProfit.Cmp(diff) < 0 {
			maxProfit = diff
		} else if maxLoss.Cmp(diff) > 0 {
			maxLoss = diff
		}
	}

	return maxProfit, maxLoss, nil
}
