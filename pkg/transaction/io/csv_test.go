package io_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/wstonks/pkg/transaction"
	"github.com/wlachs/wstonks/pkg/transaction/io"
	"testing"
)

// TestTxCsvLoader_Load is a smoke-test for a well-formatted transaction input CSV.
func TestTxCsvLoader_Load(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/smoke.csv"}
	err := loader.Load(&ctx)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(ctx.Assets))
	assert.Equal(t, 5, len(ctx.Transactions))
}

// TestTxCsvLoader_Load_Empty tests loading an empty CSV file.
func TestTxCsvLoader_Load_Empty(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/empty.csv"}
	err := loader.Load(&ctx)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(ctx.Assets))
	assert.Equal(t, 0, len(ctx.Transactions))
}

// TestTxCsvLoader_Load_Missing_File tests loading a non-existing file.
func TestTxCsvLoader_Load_Missing_File(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/###.csv"}
	err := loader.Load(&ctx)

	assert.Equal(t, fmt.Errorf("failed to open file \"../../../test/data/io/###.csv\""), err)
}

// TestTxCsvLoader_Load_No_Ts tests loading a malformed CSV file without timestamp.
func TestTxCsvLoader_Load_No_Ts(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/no_ts.csv"}
	err := loader.Load(&ctx)

	assert.NotNil(t, err)
}

// TestTxCsvLoader_Load_No_AssetId tests loading a malformed CSV file without asset ID.
func TestTxCsvLoader_Load_No_AssetId(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/no_assetId.csv"}
	err := loader.Load(&ctx)

	assert.Equal(t, fmt.Errorf("failed to parse asset ID of row [1712200000000  BUY 1.23456789 50.12345]"), err)
}

// TestTxCsvLoader_Load_No_Type tests loading a malformed CSV file without transaction type.
func TestTxCsvLoader_Load_No_Type(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/no_type.csv"}
	err := loader.Load(&ctx)

	assert.Equal(t, fmt.Errorf("failed to parse trade type of row [1712200000000 A  1.23456789 50.12345]"), err)
}

// TestTxCsvLoader_Load_No_Quantity tests loading a malformed CSV file without quantity.
func TestTxCsvLoader_Load_No_Quantity(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/no_quantity.csv"}
	err := loader.Load(&ctx)

	assert.Equal(t, fmt.Errorf("failed to parse quantity of row [1712200000000 A BUY  50.12345]"), err)
}

// TestTxCsvLoader_Load_No_UnitPrice tests loading a malformed CSV file without unit price.
func TestTxCsvLoader_Load_No_UnitPrice(t *testing.T) {
	t.Parallel()

	ctx := transaction.Context{}
	loader := io.TxCsvLoader{Path: "../../../test/data/io/transactions/no_unitPrice.csv"}
	err := loader.Load(&ctx)

	assert.Equal(t, fmt.Errorf("failed to parse unit price of row [1712200000000 A BUY 1.23456789 ]"), err)
}
