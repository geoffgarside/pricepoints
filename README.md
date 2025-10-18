# Price Points

A Go package for adjusting product prices to conform to specific price point
constraints. Price points are the units digit of a price (0-9), commonly used in
retail pricing strategies where prices end in psychologically appealing digits
(e.g., £2.99, £3.49, £3.95).

## Overview

The `pricepoints` package automatically adjusts product prices to match
predefined price points while respecting minimum and maximum price constraints.

This is useful for:

- **Psychological pricing strategies** - Prices ending in certain digits
  (e.g. 9, 5, 99) are perceived as more attractive to customers
- **Price standardization** - Ensuring all products follow consistent pricing
  rules across your catalog
- **Automated repricing** - Adjusting prices algorithmically while maintaining
  pricing rules

## Installation

```bash
go get github.com/geoffgarside/pricepoints
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/geoffgarside/pricepoints"
)

func main() {
    // Create a calculator with price points 3, 5, and 9
    // This means prices should end in .x3, .x5, or .x9
    calc, err := pricepoints.NewCalculator(3, 5, 9)
    if err != nil {
        log.Fatal(err)
    }

    // Calculate new price for a product
    // Current: £3.84, Min: £3.79, Max: £3.99
    // Prices in pence (smallest currency unit)
    newPrice, err := calc.NewPrice(384, 379, 399)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("New price: £%.2f\n", float64(newPrice)/100)
    // Output: New price: £3.85
}
```

## Usage

### Creating a Calculator

Create a calculator by specifying valid price points (0-9):

```go
calc, err := pricepoints.NewCalculator(3, 5, 9)
if err != nil {
    // Handle error (invalid price points)
}
```

Price points must be:

- Between 0 and 9 (inclusive)
- At least one price point must be provided

### Calculating New Prices

Calculate a new price that matches your price points:

```go
// All values in smallest currency unit (pence, cents, etc.)
newPrice, err := calc.NewPrice(current, min, max)
if err != nil {
    if errors.Is(err, pricepoints.ErrNoValidPrice) {
        // No valid price exists in the range
    }
}
```

**Important:** Always use the smallest currency unit (pence for GBP, cents for
USD/EUR) to avoid floating-point arithmetic issues.

### Price Preference

By default, when two prices are equidistant from the current price, the higher
price is selected (revenue maximization):

```go
calc.PreferGreaterPrices() // Default behavior
```

You can configure the calculator to prefer lower prices (customer satisfaction):

```go
calc.PreferLowerPrices()
```

**Example:**

- Current price: £3.84 (384 pence)
- Price range: £3.79 - £3.99
- Price points: [3, 5, 9]
- Valid options: £3.83 (1p lower) or £3.85 (1p higher)
- `PreferGreaterPrices()` → £3.85
- `PreferLowerPrices()` → £3.83

## Algorithm

The price calculation follows these steps:

1. **Clamp to range:** If current price is outside [min, max], adjust it to the
   nearest boundary
2. **Check existing:** If the price already matches a price point, return it
3. **Search both directions:** Find the nearest matching price point both higher
   and lower
4. **Select closest:** Return the price with the smallest difference
5. **Break ties:** If both distances are equal, use the configured preference

## Examples

### Example 1: Basic Usage

```go
calc, _ := pricepoints.NewCalculator(3, 5, 9)

// Product with current price £3.84
newPrice, _ := calc.NewPrice(384, 379, 399)
fmt.Println(newPrice) // 385 (£3.85)
```

### Example 2: Price Outside Range

```go
calc, _ := pricepoints.NewCalculator(3, 5, 9)

// Current price £5.57 is below minimum £5.82
newPrice, _ := calc.NewPrice(557, 582, 588)
fmt.Println(newPrice) // 583 (£5.83)
// Price clamped to min (582), then adjusted to nearest price point (583)
```

### Example 3: No Valid Price

```go
calc, _ := pricepoints.NewCalculator(3, 5, 9)

// Range £0.70-£0.72 contains no prices ending in 3, 5, or 9
newPrice, err := calc.NewPrice(65, 70, 72)
if errors.Is(err, pricepoints.ErrNoValidPrice) {
    fmt.Println("No valid price in range")
}
```

### Example 4: Prefer Lower Prices

```go
calc, _ := pricepoints.NewCalculator(3, 5, 9)
calc.PreferLowerPrices()

// Both £3.83 and £3.85 are valid, but lower is preferred
newPrice, _ := calc.NewPrice(384, 379, 399)
fmt.Println(newPrice) // 383 (£3.83)
```

## Demo Application

A command-line tool is included for processing CSV files of products:

```bash
cd cmd/demo
go build

# Process products.csv with price points 3, 5, 9
./demo -price-points=3,5,9 products.csv

# Prefer lower prices when equidistant
./demo -price-points=3,5,9 -prefer-lower-prices products.csv
```

### CSV Format

Input CSV should have columns:

```
Product Name,Current Price,Minimum Price,Maximum Price
Bacon,2.20,2.12,2.19
Cheese,5.57,5.82,5.88
```

Output includes the calculated new price:

```
Product Name,Current Price,Minimum Price,Maximum Price,New Price
Bacon,2.20,2.12,2.19,2.19
Cheese,5.57,5.82,5.88,5.83
```

## Error Handling

The package defines these errors:

- `ErrPricePointsMissing` - No price points provided to NewCalculator
- `ErrPricePointTooSmall` - Price point < 0
- `ErrPricePointTooLarge` - Price point > 9
- `ErrNoValidPrice` - No price in [min, max] matches any price point

Example:

```go
calc, err := pricepoints.NewCalculator(3, 5, 9)
if err != nil {
    switch {
    case errors.Is(err, pricepoints.ErrPricePointsMissing):
        // Handle missing price points
    case errors.Is(err, pricepoints.ErrPricePointTooSmall):
        // Handle invalid price point
    case errors.Is(err, pricepoints.ErrPricePointTooLarge):
        // Handle invalid price point
    }
}

newPrice, err := calc.NewPrice(current, min, max)
if errors.Is(err, pricepoints.ErrNoValidPrice) {
    // Handle case where no valid price exists
}
```

## Use Cases

### Psychological Pricing

Configure prices to end in 9 for the "charm pricing" effect:

```go
calc, _ := pricepoints.NewCalculator(9)
```

### Multiple Price Points

Use multiple price points for flexibility:

```go
// .x3, .x5, .x9 endings (e.g., £2.99, £3.49, £4.95)
calc, _ := pricepoints.NewCalculator(3, 5, 9)
```

### Even Pricing

Force all prices to even numbers:

```go
calc, _ := pricepoints.NewCalculator(0, 2, 4, 6, 8)
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Contributing

Contributions are welcome! Please ensure:

- All tests pass
- Code follows Go conventions
- New features include tests and documentation
