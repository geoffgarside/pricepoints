package pricepoints

import "testing"

func TestCalculator_PreferMethods(t *testing.T) {
	c, err := NewCalculator(3, 5, 9)
	if err != nil {
		t.Errorf("Expected NewCalculator not to return error, got %v", err)
	}

	if !c.preferGreater {
		t.Errorf("Expected Calculator to be created with preferGreater set to true, got %v",
			c.preferGreater)
	}

	c.PreferLowerPrices()

	if c.preferGreater {
		t.Errorf("Expected Calculator.PreferLowerPrices() to set preferGreater to false, got %v",
			c.preferGreater)
	}

	c.PreferGreaterPrices()

	if !c.preferGreater {
		t.Errorf("Expected Calculator.PreferGreaterPrices() to set preferGreater to true, got %v",
			c.preferGreater)
	}
}

func TestCalculator_regression1(t *testing.T) {
	c, err := NewCalculator(3, 5, 9)
	if err != nil {
		t.Errorf("Expected NewCalculator not to return error, got %v", err)
	}

	hP, hD, hE := c.nextHighestPrice(380, 399)
	if hP != 383 && hD != 3 && hE != nil {
		t.Errorf("Expected nextHighestPrice to return %v, %v, %v got %v, %v, %v",
			383, 3, nil, hP, hD, hE)
	}

	lP, lD, lE := c.nextLowestPrice(380, 379)
	if lP != 379 && lD != -1 && lE != nil {
		t.Errorf("Expected nextHighestPrice to return %v, %v, %v got %v, %v, %v",
			379, -1, nil, lP, lD, lE)
	}
}
