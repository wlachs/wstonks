package transaction

import (
	"math/big"
	"slices"
)

// GetRealizedProfit sums up the earnings for every transaction that was sold higher than the initial price.
func (ctx *Context) GetRealizedProfit() *big.Rat {
	profit := ctx.getRealizedProfitsAndLosses()

	p := big.NewRat(0, 1)
	zero := big.NewRat(0, 1)
	for _, diff := range profit {
		if zero.Cmp(diff) < 0 {
			p.Add(p, diff)
		}
	}

	return p
}

// GetRealizedLoss sums up the earnings for every transaction that was sold lower than the initial price.
func (ctx *Context) GetRealizedLoss() *big.Rat {
	profit := ctx.getRealizedProfitsAndLosses()

	p := big.NewRat(0, 1)
	zero := big.NewRat(0, 1)
	for _, diff := range profit {
		if diff.Cmp(zero) < 0 {
			p.Sub(p, diff)
		}
	}

	return p
}

// getRealizedProfitsAndLosses returns a slice of profits and losses realized with every individual SELL transaction.
func (ctx *Context) getRealizedProfitsAndLosses() []*big.Rat {
	m := map[*TxAsset][]Position{}
	var profit []*big.Rat

	for i := range ctx.Assets {
		asset := ctx.Assets[i]

		// sort transactions according to timestamp
		slices.SortFunc(asset.Transactions, func(a, b *Tx) int {
			return a.Timestamp.Compare(b.Timestamp)
		})

		for _, transaction := range asset.Transactions {
			if transaction.Type == BUY {
				m[asset] = append(m[asset], transaction.Position)
			} else if transaction.Type == SELL {
				_, diffs := subtractAssetPosition(m[asset], transaction.Position)
				profit = append(profit, diffs...)
			} else if transaction.Type == DIVIDEND {
				profit = append(profit, transaction.UnitPrice)
			}
		}
	}

	return profit
}
