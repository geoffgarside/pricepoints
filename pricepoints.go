// Package pricepoints provides functionality for adjusting product prices to
// conform to specific price point constraints. Price points are the units digit
// of a price (0-9), commonly used in retail pricing strategies where prices end
// in psychologically appealing digits (e.g., £2.99, £3.49, £3.95).
//
// The package allows configuration of valid price points and automatically
// adjusts product prices to the nearest valid price point while respecting
// minimum and maximum price constraints.
package pricepoints

import (
	"slices"
	"sort"
)

// Calculator determines new prices for products, adjusting them to conform to
// a predefined set of price points while respecting minimum and maximum price
// constraints.
//
// A price point is the units digit (0-9) of a price. For example, if the price
// points are [3, 5, 9], then valid prices would be £2.83, £2.85, £2.89, etc.
//
// When multiple valid prices are equidistant from the current price, the
// Calculator can be configured to prefer either higher or lower prices.
type Calculator struct {
	pricePoints   []int
	preferGreater bool
}

// NewCalculator creates a new Calculator with the provided price points.
//
// Price points must be integers between 0 and 9 (inclusive) representing the
// units digit that prices should end with. At least one price point must be
// provided.
//
// The Calculator is configured by default to prefer higher prices when two
// valid prices are equidistant from the current price. This can be changed
// using PreferLowerPrices.
//
// Returns an error if:
//   - No price points are provided (ErrPricePointsMissing)
//   - Any price point is less than 0 (ErrPricePointTooSmall)
//   - Any price point is greater than 9 (ErrPricePointTooLarge)
//
// Example:
//
//	calc, err := pricepoints.NewCalculator(3, 5, 9)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewCalculator(pricePoints ...int) (*Calculator, error) {
	pricePoints, err := validatePricePoints(pricePoints)
	if err != nil {
		return nil, err
	}

	c := &Calculator{
		pricePoints:   pricePoints,
		preferGreater: true,
	}

	return c, nil
}

// PreferLowerPrices configures the Calculator to prefer lower prices when
// two valid prices are equidistant from the current price.
//
// This option is useful when prioritizing customer satisfaction over
// revenue maximization.
//
// Example: If current price is £3.84, min is £3.79, max is £3.99, and price
// points are [3, 5, 9], both £3.83 and £3.85 are valid. This setting will
// select £3.83.
func (c *Calculator) PreferLowerPrices() {
	c.preferGreater = false
}

// PreferGreaterPrices configures the Calculator to prefer higher prices when
// two valid prices are equidistant from the current price.
//
// This is the default behavior and is useful when prioritizing revenue
// maximization.
//
// Example: If current price is £3.84, min is £3.79, max is £3.99, and price
// points are [3, 5, 9], both £3.83 and £3.85 are valid. This setting will
// select £3.85.
func (c *Calculator) PreferGreaterPrices() {
	c.preferGreater = true
}

// NewPrice calculates a new price that conforms to the Calculator's price points
// while respecting the provided minimum and maximum price constraints.
//
// All prices should be provided in the smallest currency unit (e.g., pence for GBP,
// cents for USD/EUR) to avoid floating-point arithmetic issues.
//
// The algorithm works as follows:
//  1. If current price is outside [min, max], it is clamped to the valid range
//  2. If the adjusted price already matches a price point, it is returned
//  3. Otherwise, search for the nearest valid price point in both directions
//  4. Return the closest valid price, with ties broken by the preference setting
//
// Parameters:
//   - current: The current price in smallest currency units
//   - min: The minimum acceptable price in smallest currency units
//   - max: The maximum acceptable price in smallest currency units
//
// Returns ErrNoValidPrice if no price within [min, max] matches any price point.
//
// Example:
//
//	calc, _ := pricepoints.NewCalculator(3, 5, 9)
//	newPrice, err := calc.NewPrice(384, 379, 399)
//	// newPrice = 385 (£3.85 if working in pence)
func (c *Calculator) NewPrice(current, min, max int) (int, error) {
	// Sanitise our value, ensure current is within min:max
	newPrice := current
	if newPrice < min {
		newPrice = min
	}

	if newPrice > max {
		newPrice = max
	}

	if c.matchesPricePoint(newPrice) {
		return newPrice, nil
	}

	highPrice, highDiff, highErr := c.nextHighestPrice(newPrice, max)
	lowPrice, lowDiff, lowErr := c.nextLowestPrice(newPrice, min)

	if highErr != nil && lowErr != nil {
		return 0, highErr
	}

	if highErr != nil {
		return lowPrice, nil
	}

	if lowErr != nil {
		return highPrice, nil
	}

	if highDiff == lowDiff {
		if c.preferGreater {
			return highPrice, nil
		}

		return lowPrice, nil
	}

	if highDiff > lowDiff {
		return lowPrice, nil
	}

	return highPrice, nil
}

// matchesPricePoint checks if the given value's units digit matches any of
// the Calculator's configured price points.
func (c *Calculator) matchesPricePoint(value int) bool {
	unit := value % 10
	return slices.Contains(c.pricePoints, unit)
}

// nextHighestPrice finds the nearest price higher than or equal to the given
// price that matches a price point, within the maximum constraint.
//
// Returns the matching price, the difference from the input price, and an
// error if no valid price exists within the range.
func (c *Calculator) nextHighestPrice(price, max int) (int, int, error) {
	for current := price; current < max; current++ {
		if c.matchesPricePoint(current) {
			return current, (current - price), nil
		}
	}

	return 0, 0, ErrNoValidPrice
}

// nextLowestPrice finds the nearest price lower than or equal to the given
// price that matches a price point, within the minimum constraint.
//
// Returns the matching price, the difference from the input price, and an
// error if no valid price exists within the range.
func (c *Calculator) nextLowestPrice(price, min int) (int, int, error) {
	for current := price; current >= min; current-- {
		if c.matchesPricePoint(current) {
			return current, (price - current), nil
		}
	}

	return 0, 0, ErrNoValidPrice
}

// validatePricePoints validates and normalizes the provided price points.
//
// Ensures that:
//   - At least one price point is provided
//   - All price points are between 0 and 9 (inclusive)
//
// Returns a sorted copy of the price points if valid, or an error otherwise.
func validatePricePoints(pricePoints []int) ([]int, error) {
	if len(pricePoints) == 0 {
		return nil, ErrPricePointsMissing
	}

	// TODO(gg): indicate which point and index is invalid
	for _, point := range pricePoints {
		if point < 0 {
			return nil, ErrPricePointTooSmall
		}

		if point > 9 {
			return nil, ErrPricePointTooLarge
		}
	}

	sort.Ints(pricePoints)
	return pricePoints, nil
}
