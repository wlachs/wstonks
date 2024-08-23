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

// ratioTestSuite contains context information for testing ratio calculation.
type ratioTestSuite struct {
	suite.Suite
	ctx *calculation.Context
}

// TestRatioTestSuite initializes and executes the test suite.
func TestRatioTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ratioTestSuite))
}

// SetupTest runs before each test case.
func (suite *ratioTestSuite) SetupTest() {
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

// TestContext_GetAssetRatio runs the smoke test.
func (suite *ratioTestSuite) TestContext_GetAssetRatio() {
	m, err := suite.ctx.GetAssetRatio(suite.ctx.AssetContext.Assets)
	keys := suite.ctx.AssetContext.GetAssetKeyMap()

	assert.NoError(suite.T(), err, "should not return error")
	assert.Equal(suite.T(), big.NewRat(1785390000000000, 1855638572748941), m[keys["A"]], "asset ratio should match")
	assert.Equal(suite.T(), big.NewRat(70248572748941, 1855638572748941), m[keys["B"]], "asset ratio should match")
}

// TestContext_GetAssetRatio runs the smoke test and makes sure that the sum of asset ratios add up to 1.
func (suite *ratioTestSuite) TestContext_GetAssetRatio_One_Sum() {
	m, err := suite.ctx.GetAssetRatio(suite.ctx.AssetContext.Assets)

	assert.NoError(suite.T(), err, "should not return error")

	sum := big.NewRat(0, 1)
	for _, rat := range m {
		sum.Add(sum, rat)
	}

	assert.Equal(suite.T(), big.NewRat(1, 1), sum, "sum should be 1")
}

// TestContext_GetAssetRatio_No_Transactions runs tests without any transactions.
func (suite *ratioTestSuite) TestContext_GetAssetRatio_No_Transactions() {
	suite.ctx.TransactionContext.Transactions = []*transaction.Tx{}
	for _, a := range suite.ctx.TransactionContext.Assets {
		a.Transactions = []*transaction.Tx{}
	}

	_, err := suite.ctx.GetAssetRatio(suite.ctx.AssetContext.Assets)

	assert.EqualError(suite.T(), err, "sum of asset worth is zero", "should throw an error if assets have no worth")
}

// TestContext_GetAssetRatio_No_Assets runs the test with no asset argument.
func (suite *ratioTestSuite) TestContext_GetAssetRatio_No_Assets() {
	_, err := suite.ctx.GetAssetRatio([]*asset.Asset{})
	assert.EqualError(suite.T(), err, "sum of asset worth is zero", "should throw an error if assets have no worth")
}

// TestContext_GetAssetRatio_Zero_Portfolio_Worth runs the test with a complete loss to make sure the correct error is returned
func (suite *ratioTestSuite) TestContext_GetAssetRatio_Zero_Portfolio_Worth() {
	for _, a := range suite.ctx.AssetContext.Assets {
		a.UnitPrice = big.NewRat(0, 1)
	}

	_, err := suite.ctx.GetAssetRatio(suite.ctx.AssetContext.Assets)

	assert.EqualError(suite.T(), err, "sum of asset worth is zero", "should throw an error if assets have no worth")
}

// TestContext_GetAssetRatio_No_Context runs the test without transaction and asset contexts.
func (suite *ratioTestSuite) TestContext_GetAssetRatio_No_Context() {
	suite.ctx.AssetContext = nil
	suite.ctx.TransactionContext = nil

	_, err := suite.ctx.GetAssetRatio([]*asset.Asset{})

	assert.EqualError(suite.T(), err, "asset or transaction context is missing", "should return error when asset / tx context is nil")
}
