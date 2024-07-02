package transaction

import (
	"fmt"
	"math/big"
	"slices"
)

// GetReturnForUnitPrice calculates the difference between the initial value of the position and its current value.
func (p Position) GetReturnForUnitPrice(unitPrice *big.Rat) *big.Rat {
	diff := big.NewRat(0, 1).Set(unitPrice)
	diff.Sub(diff, p.UnitPrice)
	diff.Mul(diff, p.Quantity)

	return diff
}

// GetAssetPositionSliceMap maps quantities to a chronologically ordered slice of positions.
// Transactions are used as a basis: There are two scenarios, BUY and SELL. In case of a BUY transaction, the position is simply added to
// the end of the position slice. In case of a SELL transaction however, the position quantity is subtracted from the oldest position.
// If the transaction value is higher than the first position, remove the first position, subtract the quantity from the transaction
// quantity and try again.
func (ctx *Context) GetAssetPositionSliceMap() map[*TxAsset][]Position {
	m := map[*TxAsset][]Position{}
	for i := range ctx.Assets {
		asset := ctx.Assets[i]

		m[asset] = ctx.GetAssetPositions(asset)

		if len(m[asset]) == 0 {
			delete(m, asset)
		}
	}

	return m
}

// GetAssetPositions calculates the open positions for the given TxAsset.
func (ctx *Context) GetAssetPositions(a *TxAsset) []Position {
	var p []Position

	// sort transactions according to timestamp
	slices.SortFunc(a.Transactions, func(a, b *Tx) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	for _, transaction := range a.Transactions {
		if transaction.Type == BUY {
			p = append(p, transaction.Position)
		} else if transaction.Type == SELL {
			p, _ = subtractAssetPosition(p, transaction.Position)
		}
	}

	return p
}

// subtractAssetPosition subtracts the position quantity from the oldest position of the asset. Returns a slice containing the profits and
// losses realized on each open position.
func subtractAssetPosition(p []Position, position Position) ([]Position, []*big.Rat) {
	if len(p) == 0 {
		return p, []*big.Rat{}
	}
	oldestPosition := p[0]

	realized := big.NewRat(0, 1).Set(position.UnitPrice)
	realized.Sub(realized, oldestPosition.UnitPrice)

	if oldestPosition.Quantity.Cmp(position.Quantity) > 0 {
		oldestPosition.Quantity.Sub(oldestPosition.Quantity, position.Quantity)

		realized.Mul(realized, position.Quantity)
		return p, []*big.Rat{realized}

	} else {
		position.Quantity.Sub(position.Quantity, oldestPosition.Quantity)
		p = p[1:]

		realized.Mul(realized, oldestPosition.Quantity)
		pp, r := subtractAssetPosition(p, position)
		return pp, append(r, realized)
	}
}

// GetAssetKeyPositions calculates the open positions for the given TxAsset key.
func (ctx *Context) GetAssetKeyPositions(assetId string) ([]Position, error) {
	for _, asset := range ctx.Assets {
		if asset.Id == assetId {
			return ctx.GetAssetPositions(asset), nil
		}
	}
	return nil, fmt.Errorf("asset with key \"%s\" not found", assetId)
}
