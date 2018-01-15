package pricepoints

import "sort"

// Calculator determines new prices of items conforming to a given
// set of price points.
type Calculator struct {
	pricePoints   []int
	preferGreater bool
}

// NewCalculator creates a new Calculator with the provided price points.
//
// The provided price points are validated to ensure they conform to the
// expected constraints.
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

// PreferLowerPrices allows for the Calculator to be configured to
// prefer lower rather than greater prices when two prices would
// conform to the new price points.
func (c *Calculator) PreferLowerPrices() {
	c.preferGreater = false
}

// PreferGreaterPrices allows for the Calculator to be configured to
// prefer higher rather than greater prices when two prices would
// conform to the new price points.
//
// This is the default behaviour.
func (c *Calculator) PreferGreaterPrices() {
	c.preferGreater = true
}

// NewPrice returns the new price of the product based on the current, min and max
// price ranges. The new price will conform to the price points of the Calculator
// or an error will be returned.
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

func (c *Calculator) matchesPricePoint(value int) bool {
	unit := value % 10
	for _, point := range c.pricePoints {
		if point == unit {
			return true
		}
	}

	return false
}

func (c *Calculator) nextHighestPrice(price, max int) (int, int, error) {
	for current := price; current < max; current++ {
		if c.matchesPricePoint(current) {
			return current, (current - price), nil
		}
	}

	return 0, 0, ErrNoValidPrice
}

func (c *Calculator) nextLowestPrice(price, min int) (int, int, error) {
	for current := price; current >= min; current-- {
		if c.matchesPricePoint(current) {
			return current, (price - current), nil
		}
	}

	return 0, 0, ErrNoValidPrice
}

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
