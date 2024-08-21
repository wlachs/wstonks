package asset

import (
	"fmt"
	"log"
	"math/big"
	"slices"
)

// Context holding Asset data
type Context struct {
	Assets []*Asset
}

// AddAssets adds a slice of Asset objects to the Context. If the asset can already be found in the Context, the price is updated to the
// newly imported value.
func (ctx *Context) AddAssets(assets []*Asset) error {
	var err error
	for _, t := range assets {
		e := ctx.addAssetInternal(t, false)
		if e != nil {
			log.Printf("failed to add asset %v: %v\n", t, e)
			err = e
		}
	}

	if err != nil {
		return err
	}

	return ctx.ValidateContext()
}

// AddAsset adds an Asset to the Context.
func (ctx *Context) AddAsset(asset *Asset) error {
	return ctx.addAssetInternal(asset, true)
}

// addAssetInternal adds the Asset to the Context.
// If the validate parameter is true, the Context will be validated after adding the Asset.
func (ctx *Context) addAssetInternal(asset *Asset, validate bool) error {
	err := updateAssets(ctx, asset)
	if err != nil {
		return err
	}

	if validate {
		return ctx.ValidateContext()
	}

	return nil
}

// updateAssets adds the Asset to the quantities in the Context.
func updateAssets(ctx *Context, asset *Asset) error {
	if asset.Id == "" {
		return fmt.Errorf("missing asset ID %v", asset)
	}

	i := slices.IndexFunc(ctx.Assets, func(a *Asset) bool {
		return a.Id == asset.Id
	})

	// If the asset is not yet known, add it
	if i == -1 {
		i = len(ctx.Assets)
		ctx.Assets = append(ctx.Assets, asset)
	}

	ctx.Assets[i] = asset
	return nil
}

// ValidateContext verifies that the Context is in a valid state.
func (ctx *Context) ValidateContext() error {
	i := slices.IndexFunc(ctx.Assets, func(a *Asset) bool {
		return a.UnitPrice.Cmp(big.NewRat(0, 1)) == -1
	})

	if i != -1 {
		asset := ctx.Assets[i]
		unitPrice, _ := asset.UnitPrice.Float32()
		return fmt.Errorf("negative asset price %s: %f < 0\n", asset.Id, unitPrice)
	}

	return nil
}
