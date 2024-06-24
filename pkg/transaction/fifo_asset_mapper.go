package transaction

import (
	"slices"
)

// GetAssetPositionSliceMap maps assets to a chronologically ordered slice of positions.
// Transactions are used as a basis: There are two scenarios, BUY and SELL. In case of a BUY transaction, the position is simply added to
// the end of the position slice. In case of a SELL transaction however, the position quantity is subtracted from the oldest position.
// If the transaction value is higher than the first position, remove the first position, subtract the quantity from the transaction
// quantity and try again.
func GetAssetPositionSliceMap(ctx Context) map[*Asset][]Position {
	m := map[*Asset][]Position{}
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
				subtractAssetPosition(m, transaction.Position)
			}
		}

		if len(m[asset]) == 0 {
			delete(m, asset)
		}
	}

	return m
}

// subtractAssetPosition subtracts the position quantity from the oldest position of the asset.
func subtractAssetPosition(m map[*Asset][]Position, position Position) {
	asset := position.Asset

	if len(m[asset]) == 0 {
		return
	}
	oldestPosition := &m[asset][0]

	if oldestPosition.Quantity.Cmp(&position.Quantity) > 0 {
		oldestPosition.Quantity.Sub(&oldestPosition.Quantity, &position.Quantity)
	} else {
		position.Quantity.Sub(&position.Quantity, &oldestPosition.Quantity)
		m[asset] = m[asset][1:]
		subtractAssetPosition(m, position)
	}
}
