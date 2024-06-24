package io

import (
	"github.com/wlachs/wstonks/pkg/transaction"
)

// TransactionLoader interface to allow populating transaction.Context with context data.
type TransactionLoader interface {
	// Load loads data to the context from an arbitrary source.
	Load(ctx *transaction.Context) error
}
