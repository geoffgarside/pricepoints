package pricepoints

import "errors"

// Package errors returned by Calculator operations.
var (
	// ErrPricePointsMissing is returned when attempting to create a Calculator
	// without providing any price points.
	ErrPricePointsMissing = errors.New("pricepoints: points not provided")

	// ErrPricePointTooSmall is returned when a price point is less than 0.
	// Price points must be between 0 and 9 (inclusive).
	ErrPricePointTooSmall = errors.New("pricepoints: point too small, less than 0")

	// ErrPricePointTooLarge is returned when a price point is greater than 9.
	// Price points must be between 0 and 9 (inclusive).
	ErrPricePointTooLarge = errors.New("pricepoints: point too large, greater than 9")

	// ErrNoValidPrice is returned by NewPrice when no price within the
	// [min, max] range matches any of the configured price points.
	// This typically occurs when the price range is too narrow to contain
	// a price matching the required price points.
	ErrNoValidPrice = errors.New("pricepoints: no valid new price for product")
)
