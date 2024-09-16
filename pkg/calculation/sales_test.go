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
	r := big.NewRat(111651, 1000)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 3, len(sales), "three assets should be sold")
	assert.Equal(suite.T(), big.NewRat(15, 1), sales[assets["A"]], "sell volume should match")
	assert.Equal(suite.T(), big.NewRat(2, 1), sales[assets["C"]], "sell volume should match")
	assert.Equal(suite.T(), big.NewRat(1, 1), sales[assets["E"]], "sell volume should match")
}

// TestGetSalesForReturn_Profit_Too_Much calculates sales required for the given profit without asset restrictions.
// This test case tries to calculate an impossible return.
func (suite *salesTestSuite) TestGetSalesForReturn_Profit_Too_Much() {
	r := big.NewRat(113652, 1000)
	_, err := suite.ctx.GetSalesForReturn(r)

	assert.EqualError(suite.T(), err, "not enough assets to sell", "should return error")
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
	r := big.NewRat(-190027653319736, 10000000000000)
	sales, err := suite.ctx.GetSalesForReturn(r)
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), 2, len(sales), "two assets should be sold")
	assert.Equal(suite.T(), big.NewRat(123456789, 100000000), sales[assets["B"]], "sell volume should match")
	assert.Equal(suite.T(), big.NewRat(2, 1), sales[assets["D"]], "sell volume should match")
}

// TestGetSalesForReturn_Loss_Too_Much calculates sales required for the given loss without asset restrictions.
// This test case tries to calculate an impossible return.
func (suite *salesTestSuite) TestGetSalesForReturn_Loss_Too_Much() {
	r := big.NewRat(-210027653319737, 10000000000000)
	_, err := suite.ctx.GetSalesForReturn(r)

	assert.EqualError(suite.T(), err, "not enough assets to sell", "should return error")
}

// TestGetMaxProfitAndLoss calculates the highest possible profit and loss based on the currently held assets.
func (suite *salesTestSuite) TestGetMaxProfitAndLoss() {
	profit, loss, err := suite.ctx.GetMaxProfitAndLoss()
	assets := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")

	profitMap := map[string]*big.Rat{
		"A": big.NewRat(106851, 1000),
		"B": big.NewRat(0, 1),
		"C": big.NewRat(8, 10),
		"D": big.NewRat(0, 1),
		"E": big.NewRat(4, 1),
		"F": big.NewRat(0, 1),
	}

	lossMap := map[string]*big.Rat{
		"A": big.NewRat(0, 1),
		"B": big.NewRat(-13753456664967, 1250000000000),
		"C": big.NewRat(0, 1),
		"D": big.NewRat(-8, 1),
		"E": big.NewRat(-2, 1),
		"F": big.NewRat(0, 1),
	}

	for a, p := range profitMap {
		assert.Equal(suite.T(), p, profit[assets[a]], "profit should match")
	}

	for a, l := range lossMap {
		assert.Equal(suite.T(), l, loss[assets[a]], "loss should match")
	}
}
