package calculation

import (
	"github.com/wlachs/wstonks/pkg/asset"
	"github.com/wlachs/wstonks/pkg/transaction"
)

// Context holding data required for live calculations.
type Context struct {
	AssetContext       *asset.Context
	TransactionContext *transaction.Context
}
