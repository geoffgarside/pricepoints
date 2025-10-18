// Demo application for the pricepoints package.
//
// This tool processes CSV files containing product pricing information and
// calculates new prices that conform to specified price points.
//
// Usage:
//
//	demo -price-points=3,5,9 products.csv
//	demo -price-points=3,5,9 -prefer-lower-prices products.csv
//
// Input CSV format (with header):
//
//	Product Name,Current Price,Minimum Price,Maximum Price
//	Bacon,2.20,2.12,2.19
//	Cheese,5.57,5.82,5.88
//
// Output includes the calculated new price:
//
//	Product Name,Original Price,Minimum Price,Maximum Price,New Price
//	Bacon,2.20,2.12,2.19,2.19
//	Cheese,5.57,5.82,5.88,5.83
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/geoffgarside/pricepoints"
)

func main() {
	var (
		pricePointsStr    string
		preferLowerPrices bool
	)

	flag.StringVar(&pricePointsStr, "price-points", "", "Comma separated list of price points")
	flag.BoolVar(&preferLowerPrices, "prefer-lower-prices", false, "Selects lower prices rather than higher prices")

	flag.Parse()

	points := strings.Split(pricePointsStr, ",")
	pricePoints := make([]int, len(points))

	for i, point := range points {
		n, err := strconv.Atoi(point)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse price points, %s: %v", point, err)
			os.Exit(1)
		}

		pricePoints[i] = n
	}

	c, err := pricepoints.NewCalculator(pricePoints...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load price points, %v", err)
		os.Exit(1)
	}

	if preferLowerPrices {
		c.PreferLowerPrices()
	}

	out := csv.NewWriter(os.Stdout)
	out.Write([]string{"Product Name", "Original Price", "Minimum Price", "Maximum Price", "New Price"})

	for _, in := range flag.Args() {
		err := updateFile(in, c, out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load products, %v", err)
			os.Exit(1)
		}
	}
}

// updateFile reads a CSV file of products, calculates new prices using the
// provided Calculator, and writes the results to the output CSV writer.
//
// The input CSV must have a header row with columns:
// Product Name, Current Price, Minimum Price, Maximum Price
//
// Each product row is processed and a new price is calculated. If no valid
// price can be determined, "no-valid-price" is written in the new price column.
func updateFile(in string, c *pricepoints.Calculator, out *csv.Writer) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}

	defer f.Close()

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		if err == io.EOF {
			return nil
		}

		return err
	}

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		currentPrice, err := decimal.NewFromString(row[1])
		if err != nil {
			return err
		}

		minimumPrice, err := decimal.NewFromString(row[2])
		if err != nil {
			return err
		}

		maximumPrice, err := decimal.NewFromString(row[3])
		if err != nil {
			return err
		}

		scale := decimal.New(100, 0)
		current := currentPrice.Mul(scale).IntPart()
		minimum := minimumPrice.Mul(scale).IntPart()
		maximum := maximumPrice.Mul(scale).IntPart()

		newPrice, err := c.NewPrice(int(current), int(minimum), int(maximum))
		if err != nil {
			if err != pricepoints.ErrNoValidPrice {
				return err
			}

			row = append(row, "no-valid-price")
		} else {
			price := decimal.New(int64(newPrice), -2)
			row = append(row, price.StringFixed(2))
		}

		out.Write(row)
		out.Flush()
	}
}
