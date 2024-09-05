package calculation_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wlachs/wstonks/pkg/asset"
	assetio "github.com/wlachs/wstonks/pkg/asset/io"
	"github.com/wlachs/wstonks/pkg/calculation"
	"github.com/wlachs/wstonks/pkg/transaction"
	txio "github.com/wlachs/wstonks/pkg/transaction/io"
	"math/big"
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
	txCsv := txio.TxCsvLoader{Path: "../../test/data/io/transactions/smoke_sales.csv"}
	err := txCsv.Load(&txCtx)

	if err != nil {
		assert.Failf(suite.T(), "failed to load transaction context: %s", err.Error())
	}

	assetCtx := asset.Context{}
	assetCsv := assetio.LiveAssetCsvLoader{Path: "../../test/data/io/assets/smoke_sales.csv"}
	err = assetCsv.Load(&assetCtx)

	if err != nil {
		assert.Failf(suite.T(), "failed to load asset context: %s", err.Error())
	}

	suite.ctx = &calculation.Context{
		AssetContext:       &assetCtx,
		TransactionContext: &txCtx,
	}
}

// TestGetSalesForReturn_Profit calculates sales required for the given profit without asset restrictions.
func (suite *salesTestSuite) TestGetSalesForReturn_Profit() {
	r := big.NewRat(1, 1)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 1, len(sales), "only one asset should be sold")
	assert.Equal(suite.T(), big.NewRat(10000, 71234), sales[assets["A"]], "sell volume should match")
}

// TestGetSalesForReturn_Profit_Complex calculates sales required for the given profit without asset restrictions.
// This test case requires selling multiple assets.
func (suite *salesTestSuite) TestGetSalesForReturn_Profit_Complex() {
	r := big.NewRat(107351, 1000)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 2, len(sales), "only one asset should be sold")
	assert.Equal(suite.T(), big.NewRat(15, 1), sales[assets["A"]], "sell volume should match")
	assert.Equal(suite.T(), big.NewRat(1, 1), sales[assets["C"]], "sell volume should match")
}

// TestGetSalesForReturn_Loss calculates sales required for the given loss without asset restrictions.
func (suite *salesTestSuite) TestGetSalesForReturn_Loss() {
	r := big.NewRat(-1, 1)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 1, len(sales), "only one asset should be sold")
	assert.Equal(suite.T(), big.NewRat(100000, 891224), sales[assets["B"]], "sell volume should match")
}

// TestGetSalesForReturn_Loss_Complex calculates sales required for the given loss without asset restrictions.
// This test case requires selling multiple assets.
func (suite *salesTestSuite) TestGetSalesForReturn_Loss_Complex() {
	r := big.NewRat(-160027653319736, 10000000000000)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 2, len(sales), "only one asset should be sold")
	assert.Equal(suite.T(), big.NewRat(123456789, 100000000), sales[assets["B"]], "sell volume should match")
	assert.Equal(suite.T(), big.NewRat(1, 1), sales[assets["D"]], "sell volume should match")
}
