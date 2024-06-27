package transaction

import (
	"math/big"
	"time"
)

// TxAsset represents an asset used in a transaction
type TxAsset struct {
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
	Asset     *TxAsset
	Timestamp time.Time
	UnitPrice *big.Rat
	Quantity  *big.Rat
}

// Tx represents a single transaction of an TxAsset.
type Tx struct {
	Position
	Type TxType
}
