package transaction

import (
	"fmt"
	"log"
	"slices"
)

// Context holding historical trade and asset data.
type Context struct {
	Transactions []*Tx
	Assets       []*TxAsset
}

// AddTransactions adds a slice of Tx objects to the Context.
// The TxAsset data is dynamically filled based on the other transactions in the Context.
func (ctx *Context) AddTransactions(transactions []Tx) error {
	var err error
	for _, t := range transactions {
		e := ctx.addTransactionInternal(t, false)
		if e != nil {
			log.Printf("failed to add transaction %v: %v\n", t, e)
			err = e
		}
	}

	if err != nil {
		return err
	}

	return ctx.Validate()
}

// AddTransaction adds a Tx to the Context.
// The TxAsset is automatically created the first time it is seen in a model.Tx.
func (ctx *Context) AddTransaction(transaction Tx) error {
	return ctx.addTransactionInternal(transaction, true)
}

// addTransactionInternal adds the transaction to the Context.
// If the validate parameter is true, the Context will be validated after adding the transaction.
func (ctx *Context) addTransactionInternal(transaction Tx, validate bool) error {
	err := updateAssets(ctx, &transaction)
	if err != nil {
		return err
	}

	err = updateTransactions(ctx, &transaction)
	if err != nil {
		return err
	}

	if validate {
		return ctx.Validate()
	}

	return nil
}

// updateAssets adds the TxAsset of the Tx object to the quantities in the Context.
func updateAssets(ctx *Context, transaction *Tx) error {
	asset := transaction.Asset
	if asset == nil || asset.Id == "" {
		return fmt.Errorf("missing asset for transaction %v", transaction)
	}

	i := slices.IndexFunc(ctx.Assets, func(a *TxAsset) bool {
		return a.Id == asset.Id
	})

	// If the asset is not yet known, add it
	if i == -1 {
		ctx.Assets = append(ctx.Assets, asset)
	}

	return nil
}

// updateTransactions adds the context to the context history of the Context.
func updateTransactions(ctx *Context, transaction *Tx) error {
	if transaction.Asset == nil {
		return fmt.Errorf("missing asset for transaction %v\n", transaction)
	}

	// update transaction asset pointer
	i := slices.IndexFunc(ctx.Assets, func(a *TxAsset) bool {
		return a.Id == transaction.Asset.Id
	})
	transaction.Asset = ctx.Assets[i]

	// add transaction to context transactions
	j := len(ctx.Transactions)
	ctx.Transactions = append(ctx.Transactions, transaction)
	ctx.Transactions[j].Asset.Transactions = append(ctx.Transactions[j].Asset.Transactions, ctx.Transactions[j])
	return nil
}
