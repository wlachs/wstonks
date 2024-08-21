package calculation

import (
	"github.com/wlachs/wstonks/pkg/asset"
	"math/big"
)

// GetAssetRatio calculates the overall worth of the given assets and returns with the provided assets as keys and a *big.Rat as value
// indicating their corresponding weight in the portfolio. The sum of values is always one, meaning the ratio is always calculated amongst
// the listed assets and non-listed assets of the portfolio are ignored.
func (ctx *Context) GetAssetRatio(assets []*asset.Asset) (map[*asset.Asset]*big.Rat, error) {
	worthMapOfAssets, err := ctx.GetAssetWorthMapOfAssets(assets)
	if err != nil {
		return nil, err
	}

	worthOfAssets, err := ctx.GetAssetWorthOfAssets(assets)
	if err != nil {
		return nil, err
	}

	m := map[*asset.Asset]*big.Rat{}
	for a, worth := range worthMapOfAssets {
		w := big.NewRat(0, 1).Set(worth)
		m[a] = w.Quo(w, worthOfAssets)
	}

	return m, nil
}
