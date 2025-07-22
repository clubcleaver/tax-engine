package main

import "testing"

func TestCalculateWithInclusiveBase(t *testing.T) {

	simple10PercentTemplate := TaxApplicationTemplate{
		Applications: []Application{
			{AppType: TaxTypePercent, Percent: 10.0, ApplyOn: TaxAppOnRefVal},
		},
	}

	complexTemplate := TaxApplicationTemplate{
		Applications: []Application{
			{AppType: TaxTypePercent, Percent: 10.0, ApplyOn: TaxAppOnRefVal},
			{AppType: TaxTypeFixed, FixedAmount: 500},
			{AppType: TaxTypePercent, Percent: 5.0, ApplyOn: TaxAppOnRunningGross},
		},
	}

	// Test Cases
	testCases := []struct {
		name              string
		inclusiveBase     int
		template          TaxApplicationTemplate
		expectedExclusive int
		expectedTax       int
		expectError       bool
	}{
		{
			name:              "Simple 10 percent tax",
			inclusiveBase:     11000, // $110.00
			template:          simple10PercentTemplate,
			expectedExclusive: 10000, // $100.00
			expectedTax:       1000,  // $10.00
			expectError:       false,
		},
		{
			name:              "Complex template with fixed and cascading tax",
			inclusiveBase:     12075, // Base of 10000 -> tax is 1575 -> total 11575
			template:          complexTemplate,
			expectedExclusive: 10000,
			expectedTax:       2075,
			expectError:       false,
		},
		{
			name:          "Case requiring rounding adjustment",
			inclusiveBase: 10000, // $100.00 inclusive of 7.5% tax
			template: TaxApplicationTemplate{
				Applications: []Application{
					{AppType: TaxTypePercent, Percent: 7.5, ApplyOn: TaxAppOnRefVal},
				},
			},
			expectedExclusive: 9302,
			expectedTax:       698,
			expectError:       false,
		},
		{
			name:              "Zero inclusive base should error",
			inclusiveBase:     0,
			template:          simple10PercentTemplate,
			expectedExclusive: 0,
			expectedTax:       0,
			expectError:       true,
		},
		{
			name:              "Negative inclusive base should error",
			inclusiveBase:     -500,
			template:          simple10PercentTemplate,
			expectedExclusive: 0,
			expectedTax:       0,
			expectError:       true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			exclusive, tax, err := CalculateWithInclusiveBase(tc.inclusiveBase, &tc.template)

			// if error expected
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			// if unexpected error
			if !tc.expectError && err != nil {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			// calculated values
			if exclusive != tc.expectedExclusive {
				t.Errorf("Mismatched exclusive base: expected %d, got %d", tc.expectedExclusive, exclusive)
			}

			if tax != tc.expectedTax {
				t.Errorf("Mismatched tax: expected %d, got %d", tc.expectedTax, tax)
			}

		})

	}

}
