package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/transaction"
	"math/big"
	"slices"
	"sort"
)

// GetSalesForReturn calculates how much and which positions should be sold in order to realize the given return.
func (ctx *Context) GetSalesForReturn(r *big.Rat) (map[*asset.Asset]*big.Rat, error) {
	if r == nil {
		return nil, fmt.Errorf("return shouldn't be nil")
	}

	assetCtx := ctx.AssetContext
	if assetCtx == nil {
		return nil, fmt.Errorf("asset context not set")
	}

	return ctx.GetSalesForReturnForAssets(r, ctx.AssetContext.Assets)
}

// GetSalesForReturnForAssets calculates how much and which positions should be sold in order to realize the given return.
func (ctx *Context) GetSalesForReturnForAssets(r *big.Rat, assets []*asset.Asset) (map[*asset.Asset]*big.Rat, error) {
	if r == nil {
		return nil, fmt.Errorf("return shouldn't be nil")
	}

	profits, losses, err := ctx.GetMaxProfitAndLossForAssets(assets)
	if err != nil {
		return nil, err
	}

	zero := big.NewRat(0, 1)
	rr := big.NewRat(0, 1).Set(r)

	if r.Cmp(zero) < 0 {
		return ctx.getSalesForLossWithAssets(rr, assets, losses)
	}

	return ctx.getSalesForProfitWithAssets(rr, assets, profits)
}

// getSalesForProfitWithAssets calculates how much of the given assets have to be sold to get the given profit
func (ctx *Context) getSalesForProfitWithAssets(r *big.Rat, assets []*asset.Asset, profits map[*asset.Asset]*big.Rat) (map[*asset.Asset]*big.Rat, error) {
	sort.Slice(assets, func(i, j int) bool {
		return profits[assets[i]].Cmp(profits[assets[j]]) > 0
	})

	return ctx.sellForProfit(r, assets, profits)
}

// sellForProfit recursively iterates over the open asset positions and sells them until the desired profit is realized.
func (ctx *Context) sellForProfit(r *big.Rat, assets []*asset.Asset, profits map[*asset.Asset]*big.Rat) (map[*asset.Asset]*big.Rat, error) {
	if len(assets) == 0 {
		return nil, fmt.Errorf("not enough assets to sell")
	}

	a := assets[0]
	assets = assets[1:]
	positions, err := ctx.TransactionContext.GetAssetKeyPositions(a.Id)
	if err != nil {
		return nil, err
	}

	maxProfit := profits[a]
	m := map[*asset.Asset]*big.Rat{}
	m[a] = big.NewRat(0, 1)
	diff := big.NewRat(0, 1)

	for _, position := range positions {
		ret := position.GetReturnForUnitPrice(a.UnitPrice)
		d := big.NewRat(0, 1).Add(diff, ret)

		if d.Cmp(r) >= 0 {
			// (r - diff) / ret
			tmp := big.NewRat(0, 1).Sub(r, diff)
			tmp.Quo(tmp, ret)
			tmp.Mul(tmp, position.Quantity)

			m[a].Add(m[a], tmp)
			return m, nil

		}

		diff.Add(diff, ret)
		m[a].Add(m[a], position.Quantity)

		if d.Cmp(maxProfit) != 0 {
			continue
		}

		r.Sub(r, d)

		rest, e := ctx.sellForProfit(r, assets, profits)
		if e != nil {
			return nil, e
		}

		for k, v := range rest {
			m[k] = v
		}

		return m, nil
	}

	return nil, fmt.Errorf("not enough positions to realize the profit")
}

// getSalesForLossWithAssets calculates how much of the given assets have to be sold to get the given loss
func (ctx *Context) getSalesForLossWithAssets(r *big.Rat, assets []*asset.Asset, losses map[*asset.Asset]*big.Rat) (map[*asset.Asset]*big.Rat, error) {
	sort.Slice(assets, func(i, j int) bool {
		return losses[assets[i]].Cmp(losses[assets[j]]) < 0
	})

	return ctx.sellForLoss(r, assets, losses)
}

// sellForLoss recursively iterates over the open asset positions and sells them until the desired loss is realized.
func (ctx *Context) sellForLoss(r *big.Rat, assets []*asset.Asset, losses map[*asset.Asset]*big.Rat) (map[*asset.Asset]*big.Rat, error) {
	if len(assets) == 0 {
		return nil, fmt.Errorf("not enough assets to sell")
	}

	a := assets[0]
	assets = assets[1:]
	positions, err := ctx.TransactionContext.GetAssetKeyPositions(a.Id)
	if err != nil {
		return nil, err
	}

	maxLoss := losses[a]
	m := map[*asset.Asset]*big.Rat{}
	m[a] = big.NewRat(0, 1)
	diff := big.NewRat(0, 1)

	for _, position := range positions {
		ret := position.GetReturnForUnitPrice(a.UnitPrice)
		d := big.NewRat(0, 1).Add(diff, ret)

		if d.Cmp(r) <= 0 {
			// (r - diff) / ret
			tmp := big.NewRat(0, 1).Sub(r, diff)
			tmp.Quo(tmp, ret)
			tmp.Mul(tmp, position.Quantity)

			m[a].Add(m[a], tmp)
			return m, nil
		}

		diff.Add(diff, ret)
		m[a].Add(m[a], position.Quantity)

		if d.Cmp(maxLoss) != 0 {
			continue
		}

		r.Sub(r, d)

		rest, e := ctx.sellForLoss(r, assets, losses)
		if e != nil {
			return nil, e
		}

		for k, v := range rest {
			m[k] = v
		}

		return m, nil
	}

	return nil, fmt.Errorf("not enough positions to realize the losses")
}

// GetMaxProfitAndLoss calculates the maximum realizable profit and loss for each asset with live data.
func (ctx *Context) GetMaxProfitAndLoss() (map[*asset.Asset]*big.Rat, map[*asset.Asset]*big.Rat, error) {
	assetCtx := ctx.AssetContext
	if assetCtx == nil {
		return nil, nil, fmt.Errorf("asset context not set")
	}

	return ctx.GetMaxProfitAndLossForAssets(assetCtx.Assets)
}

// GetMaxProfitAndLossForAssets calculates the maximum realizable profit and loss for each asset with live data.
func (ctx *Context) GetMaxProfitAndLossForAssets(assets []*asset.Asset) (map[*asset.Asset]*big.Rat, map[*asset.Asset]*big.Rat, error) {
	profit := map[*asset.Asset]*big.Rat{}
	loss := map[*asset.Asset]*big.Rat{}

	for i := range assets {
		p, l, err := ctx.GetMaxProfitAndLossForAsset(assets[i])
		if err != nil {
			return nil, nil, err
		}

		profit[assets[i]] = p
		loss[assets[i]] = l
	}

	return profit, loss, nil
}

// GetMaxProfitAndLossForAsset checks every open position of the asset and calculates the maximum realizable profit and loss.
func (ctx *Context) GetMaxProfitAndLossForAsset(a *asset.Asset) (*big.Rat, *big.Rat, error) {
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
