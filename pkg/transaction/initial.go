package transaction

import (
	"math/big"
)

// GetAssetInitialWorthMap calculates the initial worth of every asset
func (ctx *Context) GetAssetInitialWorthMap() map[*TxAsset]*big.Rat {
	m := map[*TxAsset]*big.Rat{}

	for _, a := range ctx.Assets {
		m[a] = ctx.GetAssetInitialWorth(a)
	}

	return m
}

// GetAssetInitialWorth calculates the initial worth of the given asset
func (ctx *Context) GetAssetInitialWorth(a *TxAsset) *big.Rat {
	positions := ctx.GetAssetPositions(a)
	worth := big.NewRat(0, 1)

	for _, p := range positions {
		w := big.NewRat(0, 1)
		w.Mul(p.Quantity, p.UnitPrice)
		worth.Add(worth, w)
	}

	return worth
}

// GetAssetKeyInitialWorthMap calculates the initial worth of every asset
func (ctx *Context) GetAssetKeyInitialWorthMap() map[string]*big.Rat {
	m := map[string]*big.Rat{}

	for a, w := range ctx.GetAssetInitialWorthMap() {
		m[a.Id] = w
	}

	return m
}
