package io

import (
	"github.com/wlachs/wstonks/pkg/asset"
)

// LiveAssetLoader interface to allow populating asset.Context with context data.
type LiveAssetLoader interface {
	// Load loads data to the context from an arbitrary source.
	Load(ctx *asset.Context) error
}
