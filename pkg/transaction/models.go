package transaction

import (
	"math/big"
	"time"
)

// Asset represents an asset used in a transaction
type Asset struct {
	Id           string
	Transactions []*Tx
}

// TxType holds the different transaction types as a pseudo-enum.
type TxType = int

const (
	BUY TxType = iota
	SELL
	DIVIDEND
)

// Position depicts a certain quantity of an asset at a given time at a given unit price.
type Position struct {
	Asset     *Asset
	Timestamp time.Time
	UnitPrice big.Rat
	Quantity  big.Rat
}

// Tx represents a single transaction of an Asset.
type Tx struct {
	Position
	Type TxType
}
