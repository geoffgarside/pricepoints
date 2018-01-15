package pricepoints_test

import (
	"testing"

	"github.com/geoffgarside/pricepoints"
)

func TestNewCalculatorPricePointValidations(t *testing.T) {
	tests := []struct {
		PricePoints []int
		Error       error
	}{
		{[]int{}, pricepoints.ErrPricePointsMissing},
		{[]int{-1}, pricepoints.ErrPricePointTooSmall},
		{[]int{10}, pricepoints.ErrPricePointTooLarge},
	}

	for i, test := range tests {
		c, err := pricepoints.NewCalculator(test.PricePoints...)
		if err != test.Error {
			t.Errorf("#%d Expected NewCalculator to return err %v, got %v",
				i, test.Error, err)
		}

		if c != nil {
			t.Errorf("#%d Expected NewCalculator *Calculator to be nil", i)
		}
	}
}

func TestNewCalculator(t *testing.T) {
	c, err := pricepoints.NewCalculator(3, 5, 9)
	if err != nil {
		t.Errorf("Expected NewCalculator to return nil error, got %v", err)
	}

	if c == nil {
		t.Errorf("Expected NewCalculator to return non nil *Calculator, got %v", c)
	}
}

func TestCalculatorNewPrice(t *testing.T) {
	tests := []struct {
		PricePoints  []int
		ProductName  string
		CurrentPence int
		MinimumPence int
		MaximumPence int
		NewPence     int
		Error        error
	}{
		{[]int{3, 5, 9}, "Anchovies", 384, 384, 389, 385, nil},
		{[]int{3, 5, 9}, "Bacon", 220, 212, 219, 219, nil},
		{[]int{3, 5, 9}, "Cheese", 557, 582, 588, 583, nil},
		{[]int{3, 5, 9}, "Dates", 109, 88, 91, 89, nil},
		{[]int{3, 5, 9}, "Eggs", 65, 70, 72, 0, pricepoints.ErrNoValidPrice},
		{[]int{3, 5, 9}, "Fish", 384, 379, 399, 385, nil},
		{[]int{3, 5, 9}, "Ham", 77, 70, 72, 0, pricepoints.ErrNoValidPrice},
		{[]int{2, 5, 9}, "Chips", 384, 379, 399, 385, nil},
		{[]int{3, 6, 9}, "Buns", 384, 379, 399, 383, nil},
		{[]int{3, 6, 9}, "Lettuce", 383, 379, 399, 383, nil},
		{[]int{3, 6, 9}, "Tomatoes", 380, 379, 399, 379, nil},
	}

	for i, test := range tests {
		i := i + 1
		c, err := pricepoints.NewCalculator(test.PricePoints...)
		if err != nil {
			t.Fatalf("#%d Expected NewCalculator not to return error, got %v", i, err)
		}

		newPrice, err := c.NewPrice(test.CurrentPence, test.MinimumPence, test.MaximumPence)
		if err != test.Error {
			t.Errorf("#%d Expected NewPrice to return %v, got %v", i, test.Error, err)
		}

		if newPrice != test.NewPence {
			t.Errorf("#%d Expected NewPrice to return new value of %d, got %d", i, test.NewPence, newPrice)
		}
	}
}

func TestCalculatorNewPricePreferLowerPrice(t *testing.T) {
	tests := []struct {
		PricePoints  []int
		ProductName  string
		CurrentPence int
		MinimumPence int
		MaximumPence int
		NewPence     int
		Error        error
	}{
		{[]int{3, 5, 9}, "Anchovies", 384, 384, 389, 385, nil},
		{[]int{3, 5, 9}, "Bacon", 220, 212, 219, 219, nil},
		{[]int{3, 5, 9}, "Cheese", 557, 582, 588, 583, nil},
		{[]int{3, 5, 9}, "Dates", 109, 88, 91, 89, nil},
		{[]int{3, 5, 9}, "Eggs", 65, 70, 72, 0, pricepoints.ErrNoValidPrice},
		{[]int{3, 5, 9}, "Fish", 384, 379, 399, 383, nil},
		{[]int{3, 5, 9}, "Ham", 77, 70, 72, 0, pricepoints.ErrNoValidPrice},
		{[]int{2, 5, 9}, "Chips", 384, 379, 399, 385, nil},
		{[]int{3, 6, 9}, "Buns", 384, 379, 399, 383, nil},
	}

	for i, test := range tests {
		i := i + 1
		c, err := pricepoints.NewCalculator(test.PricePoints...)
		if err != nil {
			t.Fatalf("#%d Expected NewCalculator not to return error, got %v", i, err)
		}

		c.PreferLowerPrices()

		newPrice, err := c.NewPrice(test.CurrentPence, test.MinimumPence, test.MaximumPence)
		if err != test.Error {
			t.Errorf("#%d Expected NewPrice to return %v, got %v", i, test.Error, err)
		}

		if newPrice != test.NewPence {
			t.Errorf("#%d Expected NewPrice to return new value of %d, got %d", i, test.NewPence, newPrice)
		}
	}
}
