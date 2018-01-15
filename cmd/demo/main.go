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
