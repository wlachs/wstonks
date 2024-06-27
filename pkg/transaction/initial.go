package transaction

import (
	"math/big"
)

// GetAssetInitialWorthMap calculates the initial worth of every asset
func (ctx *Context) GetAssetInitialWorthMap() map[*TxAsset]*big.Rat {
	positions := ctx.GetAssetPositionSliceMap()
	m := map[*TxAsset]*big.Rat{}

	for a, pos := range positions {
		worth := big.NewRat(0, 1)

		for _, p := range pos {
			w := big.NewRat(0, 1).Set(p.Quantity)
			w.Mul(w, p.UnitPrice)
			worth.Add(worth, w)
		}

		m[a] = worth
	}

	return m
}

// GetAssetKeyInitialWorthMap calculates the initial worth of every asset
func (ctx *Context) GetAssetKeyInitialWorthMap() map[string]*big.Rat {
	m := map[string]*big.Rat{}

	for a, w := range ctx.GetAssetInitialWorthMap() {
		m[a.Id] = w
	}

	return m
}
