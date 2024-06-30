package calculation

import (
	"fmt"
	"github.com/wlachs/wstonks/pkg/asset"
	"math/big"
)

// GetDistributionAdjustmentMapWithBudget calculates how many quantities have to be bough with the given budget in order to reach the given distribution.
func (ctx *Context) GetDistributionAdjustmentMapWithBudget(distribution map[*asset.Asset]*big.Rat, budget *big.Rat) (map[*asset.Asset]*big.Rat, error) {
	err := validateDistribution(distribution)
	if err != nil {
		return nil, err
	}

	worthMap, err := ctx.GetAssetWorthMap()
	if err != nil {
		return nil, err
	}

	globalWorth := big.NewRat(0, 1).Set(budget)
	for a := range distribution {
		w, ok := worthMap[a]
		if ok {
			globalWorth.Add(globalWorth, w)
		}
	}

	m := map[*asset.Asset]*big.Rat{}
	zero := big.NewRat(0, 1)
	for a, q := range distribution {
		w, ok := worthMap[a]
		if !ok {
			w = big.NewRat(0, 1)
		}

		r := big.NewRat(0, 1).Set(globalWorth)
		r.Mul(r, q)
		r.Sub(r, w)

		if r.Cmp(zero) != 0 {
			m[a] = r
		}
	}

	return m, nil
}

// GetDistributionAdjustmentMapWithoutSelling calculates how many quantities have to be bough in order to reach the given distribution.
func (ctx *Context) GetDistributionAdjustmentMapWithoutSelling(distribution map[*asset.Asset]*big.Rat) (map[*asset.Asset]*big.Rat, error) {
	err := validateDistribution(distribution)
	if err != nil {
		return nil, err
	}

	txCtx := ctx.TransactionContext
	if txCtx == nil {
		return nil, fmt.Errorf("transaction context not set")
	}

	worthMap, err := ctx.GetAssetWorthMap()
	if err != nil {
		return nil, err
	}

	globalWorth := big.NewRat(0, 1)
	for a := range distribution {
		w, ok := worthMap[a]
		if ok {
			globalWorth.Add(globalWorth, w)
		}
	}

	var bestPerformer *asset.Asset
	var bestPerformance *big.Rat
	for a, d := range distribution {
		worthDifference := big.NewRat(0, 1)
		w, ok := worthMap[a]
		if ok {
			worthDifference.Set(w)
		}

		idealWorth := big.NewRat(0, 1).Set(globalWorth)
		idealWorth.Mul(idealWorth, d)
		worthDifference.Sub(worthDifference, idealWorth)

		if bestPerformance == nil || bestPerformance.Cmp(worthDifference) < 0 {
			bestPerformance = worthDifference
			bestPerformer = a
		}
	}

	budget, ok := worthMap[bestPerformer]
	if !ok {
		budget = big.NewRat(0, 1)
	}

	budget.Quo(budget, distribution[bestPerformer])
	budget.Sub(budget, globalWorth)

	withBudget, err := ctx.GetDistributionAdjustmentMapWithBudget(distribution, budget)
	if err != nil {
		return nil, err
	}

	return withBudget, nil
}

// validateDistribution makes sure that the overall distribution sum is not more than 1.
func validateDistribution(distribution map[*asset.Asset]*big.Rat) error {
	sum := big.NewRat(0, 1)

	for _, d := range distribution {
		sum.Add(sum, d)
	}

	if sum.Cmp(big.NewRat(1, 1)) != 0 {
		s, _ := sum.Float32()
		return fmt.Errorf("overall sum of distributed values %f ≠ 1", s)
	}

	return nil
}
