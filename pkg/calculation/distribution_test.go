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

// distributionTestSuite contains context information for testing distribution calculation.
type distributionTestSuite struct {
	suite.Suite
	ctx *calculation.Context
}

// TestDistributionTestSuite initializes and executes the test suite.
func TestDistributionTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(distributionTestSuite))
}

// SetupTest runs before each test case.
func (suite *distributionTestSuite) SetupTest() {
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

// TestGetDistributionAdjustmentMapWithBudget tests calculating how much of the individual assets to buy / sell to achieve the
// desired asset distribution. The budget parameter measures the newly introduced capital to the system.
func (suite *distributionTestSuite) TestGetDistributionAdjustmentMapWithBudget() {
	assets := suite.ctx.AssetContext.GetAssetKeyMap()
	dist := map[*asset.Asset]*big.Rat{
		assets["A"]: big.NewRat(2, 3),
		assets["B"]: big.NewRat(1, 3),
	}
	budget := big.NewRat(1000, 1)

	res, err := suite.ctx.GetDistributionAdjustmentMapWithBudget(dist, budget)

	assert.NoError(suite.T(), err, "should return no error")
	assert.Equal(suite.T(), big.NewRat(2597982154740469, 15000000000000), res[assets["A"]], "should match calculated value")
	assert.Equal(suite.T(), big.NewRat(12402017845259531, 15000000000000), res[assets["B"]], "should match calculated value")
}

// TestGetDistributionAdjustmentMapWithoutSelling tests calculating how much of the individual assets to buy to achieve the
// desired asset distribution. There is no budget set, the necessary amount is determined in the calculation.
func (suite *distributionTestSuite) TestGetDistributionAdjustmentMapWithoutSelling() {
	assets := suite.ctx.AssetContext.GetAssetKeyMap()
	dist := map[*asset.Asset]*big.Rat{
		assets["A"]: big.NewRat(2, 3),
		assets["B"]: big.NewRat(1, 3),
	}

	res, err := suite.ctx.GetDistributionAdjustmentMapWithoutSelling(dist)

	assert.NoError(suite.T(), err, "should return no error")
	assert.Nil(suite.T(), res[assets["A"]], "should be nil")
	assert.Equal(suite.T(), big.NewRat(7402017845259531, 10000000000000), res[assets["B"]], "should match calculated value")
}
