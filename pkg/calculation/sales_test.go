package calculation_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wlachs/wstonks/pkg/asset"
	assetio "github.com/wlachs/wstonks/pkg/asset/io"
	"github.com/wlachs/wstonks/pkg/calculation"
	"github.com/wlachs/wstonks/pkg/transaction"
	txio "github.com/wlachs/wstonks/pkg/transaction/io"
	"testing"
)

// salesTestSuite contains context information for testing sales calculation.
type salesTestSuite struct {
	suite.Suite
	ctx *calculation.Context
}

// TestSalesTestSuite initializes and executes the test suite.
func TestSalesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(salesTestSuite))
}

// SetupTest runs before each test case.
func (suite *salesTestSuite) SetupTest() {
	txCtx := transaction.Context{}
	txCsv := txio.TxCsvLoader{Path: "../../test/data/io/transactions/smoke.csv"}
	err := txCsv.Load(&txCtx)

	if err != nil {
		assert.Failf(suite.T(), "failed to load transaction context: %s", err.Error())
	}

	assetCtx := asset.Context{}
	assetCsv := assetio.LiveAssetCsvLoader{Path: "../../test/data/io/assets/smoke.csv"}
	err = assetCsv.Load(&assetCtx)

	if err != nil {
		assert.Failf(suite.T(), "failed to load asset context: %s", err.Error())
	}

	suite.ctx = &calculation.Context{
		AssetContext:       &assetCtx,
		TransactionContext: &txCtx,
	}
}
