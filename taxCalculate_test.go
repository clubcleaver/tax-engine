package main

import "testing"

func TestTaxCalculate(t *testing.T) {

	testCases := []struct {
		name        string // A descriptive name for the test case
		baseAmount  int
		template    TaxApplicationTemplate
		expectedTax int
		expectError bool // Should this case produce an error?
	}{
		// 2. Create a slice of test cases
		{
			name:       "Simple 10 percent tax on reference value",
			baseAmount: 10000, // $100.00
			template: TaxApplicationTemplate{
				Applications: []Application{
					{AppType: TaxTypePercent, Percent: 10.0, ApplyOn: TaxAppOnRefVal},
				},
			},
			expectedTax: 1000, // $10.00
			expectError: false,
		},
		{
			name:       "Simple fixed amount tax",
			baseAmount: 10000,
			template: TaxApplicationTemplate{
				Applications: []Application{
					{AppType: TaxTypeFixed, FixedAmount: 2500},
				},
			},
			expectedTax: 2500, // $25.00
			expectError: false,
		},
		{
			name:       "Cascading tax on running gross",
			baseAmount: 10000,
			template: TaxApplicationTemplate{
				Applications: []Application{
					{AppType: TaxTypePercent, Percent: 10.0, ApplyOn: TaxAppOnRefVal},      // Tax is 1000, running gross is 11000
					{AppType: TaxTypePercent, Percent: 5.0, ApplyOn: TaxAppOnRunningGross}, // Tax is 550 (5% of 11000)
				},
			},
			expectedTax: 1550, // Total tax is 1000 + 550
			expectError: false,
		},
		{
			name:       "Your original complex scenario",
			baseAmount: 10000,
			template: TaxApplicationTemplate{
				Applications: []Application{
					{AppType: TaxTypePercent, Percent: 10, ApplyOn: TaxAppOnRefVal},
					{AppType: TaxTypePercent, Percent: 10, ApplyOn: TaxAppOnPrevTax},
					{AppType: TaxTypePercent, Percent: 10, ApplyOn: TaxAppOnRunningGross},
					{AppType: TaxTypeFixed, FixedAmount: 10000},
				},
			},
			expectedTax: 12210,
			expectError: false,
		},
		{
			name:        "Error on negative base amount",
			baseAmount:  -100,
			template:    TaxApplicationTemplate{},
			expectedTax: 0,
			expectError: true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			tax, err := TaxCaculate(tc.baseAmount, &tc.template)

			if tc.expectError {
				if err == nil {
					t.Errorf("on TaxCalculate: Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("on TaxCalculate: Did not expect an error, but got: %v", err)
				}
			}

			if tax != tc.expectedTax {
				t.Errorf("on TaxCalculate: Expected tax to be %d, but got %d", tc.expectedTax, tax)
			}

		})

	}

}
