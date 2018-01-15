package pricepoints

import "errors"

// Errors
var (
	ErrPricePointsMissing = errors.New("pricepoints: points not provided")
	ErrPricePointTooSmall = errors.New("pricepoints: point too small, less than 0")
	ErrPricePointTooLarge = errors.New("pricepoints: point too large, greater than 9")
	ErrNoValidPrice       = errors.New("pricepoints: no valid new price for product")
)
