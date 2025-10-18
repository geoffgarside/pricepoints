package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/geoffgarside/pricepoints"
)

func TestUpdateFile(t *testing.T) {
	tests := []struct {
		name            string
		inputFile       string
		pricePoints     []int
		preferLower     bool
		expectedRows    [][]string
		expectError     bool
		errorContains   string
	}{
		{
			name:        "valid products with all valid prices",
			inputFile:   "testdata/valid_products.csv",
			pricePoints: []int{3, 5, 9},
			preferLower: false,
			expectedRows: [][]string{
				{"Bacon", "2.20", "2.12", "2.19", "2.19"},
				{"Cheese", "5.57", "5.82", "5.88", "5.83"},
				{"Fish", "3.84", "3.79", "3.99", "3.85"},
				{"Anchovies", "3.84", "3.84", "3.89", "3.85"},
				{"Dates", "1.09", "0.88", "0.91", "0.89"},
			},
			expectError: false,
		},
		{
			name:        "products with no valid prices",
			inputFile:   "testdata/products_with_no_valid_price.csv",
			pricePoints: []int{3, 5, 9},
			preferLower: false,
			expectedRows: [][]string{
				{"Eggs", "0.65", "0.70", "0.72", "no-valid-price"},
				{"Ham", "0.77", "0.70", "0.72", "no-valid-price"},
			},
			expectError: false,
		},
		{
			name:        "prefer lower prices",
			inputFile:   "testdata/mixed_products.csv",
			pricePoints: []int{3, 5, 9},
			preferLower: true,
			expectedRows: [][]string{
				{"Bacon", "2.20", "2.12", "2.19", "2.19"},
				{"Eggs", "0.65", "0.70", "0.72", "no-valid-price"},
				{"Fish", "3.84", "3.79", "3.99", "3.83"},
			},
			expectError: false,
		},
		{
			name:         "empty file after header",
			inputFile:    "testdata/empty_after_header.csv",
			pricePoints:  []int{3, 5, 9},
			preferLower:  false,
			expectedRows: [][]string{},
			expectError:  false,
		},
		{
			name:          "non-existent file",
			inputFile:     "testdata/does_not_exist.csv",
			pricePoints:   []int{3, 5, 9},
			preferLower:   false,
			expectedRows:  nil,
			expectError:   true,
			errorContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create calculator
			calc, err := pricepoints.NewCalculator(tt.pricePoints...)
			if err != nil {
				t.Fatalf("Failed to create calculator: %v", err)
			}

			if tt.preferLower {
				calc.PreferLowerPrices()
			}

			// Create output buffer
			var buf bytes.Buffer
			writer := csv.NewWriter(&buf)

			// Run updateFile
			err = updateFile(tt.inputFile, calc, writer)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.errorContains)
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Flush and read output
			writer.Flush()
			if writer.Error() != nil {
				t.Fatalf("CSV writer error: %v", writer.Error())
			}

			// Parse output
			reader := csv.NewReader(&buf)
			outputRows, err := reader.ReadAll()
			if err != nil {
				t.Fatalf("Failed to read output CSV: %v", err)
			}

			// Compare rows
			if len(outputRows) != len(tt.expectedRows) {
				t.Errorf("Expected %d rows, got %d", len(tt.expectedRows), len(outputRows))
				t.Logf("Output:\n%v", outputRows)
			}

			for i, expectedRow := range tt.expectedRows {
				if i >= len(outputRows) {
					t.Errorf("Missing row %d: expected %v", i, expectedRow)
					continue
				}

				if len(outputRows[i]) != len(expectedRow) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(expectedRow), len(outputRows[i]))
					continue
				}

				for j, expectedCol := range expectedRow {
					if outputRows[i][j] != expectedCol {
						t.Errorf("Row %d, Col %d: expected %q, got %q", i, j, expectedCol, outputRows[i][j])
					}
				}
			}
		})
	}
}

func TestUpdateFileInvalidPriceFormat(t *testing.T) {
	// Create a temporary file with invalid price format
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid_prices.csv")

	content := `Product Name,Current Price,Minimum Price,Maximum Price
InvalidProduct,not-a-number,2.12,2.19
`
	err := os.WriteFile(invalidFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	calc, err := pricepoints.NewCalculator(3, 5, 9)
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	err = updateFile(invalidFile, calc, writer)
	if err == nil {
		t.Error("Expected error for invalid price format, got nil")
	}
}

func TestUpdateFileInvalidMinPrice(t *testing.T) {
	// Create a temporary file with invalid minimum price format
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid_min.csv")

	content := `Product Name,Current Price,Minimum Price,Maximum Price
Product,2.20,invalid,2.19
`
	err := os.WriteFile(invalidFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	calc, err := pricepoints.NewCalculator(3, 5, 9)
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	err = updateFile(invalidFile, calc, writer)
	if err == nil {
		t.Error("Expected error for invalid minimum price format, got nil")
	}
}

func TestUpdateFileInvalidMaxPrice(t *testing.T) {
	// Create a temporary file with invalid maximum price format
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid_max.csv")

	content := `Product Name,Current Price,Minimum Price,Maximum Price
Product,2.20,2.12,invalid
`
	err := os.WriteFile(invalidFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	calc, err := pricepoints.NewCalculator(3, 5, 9)
	if err != nil {
		t.Fatalf("Failed to create calculator: %v", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	err = updateFile(invalidFile, calc, writer)
	if err == nil {
		t.Error("Expected error for invalid maximum price format, got nil")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
