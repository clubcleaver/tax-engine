# Overview
It is designed to handle requirements such as:

* Applying multiple taxes sequentially.
* Calculating taxes based on a running total (cascading taxes).
* Solving the difficult problem of finding the pre-tax base from a tax-inclusive total.

This project is a showcase of how to model and solve a complex business problem in Go, with a focus on correctness, financial precision, and robust 


## Key Features
* Rules-Based Engine: Define complex tax structures using declarative TaxApplicationTemplate structs.
* Complex Templates: Supports percentage-based and fixed-amount taxes applied sequentially.
* Multiple Calculation Bases: Taxes can be calculated on the initial reference value, the previously calculated tax amount, or the running gross total.
* Inclusive Price Calculation: Implements a robust iterative algorithm to accurately determine the exclusive base and tax amounts from a tax-inclusive total.
* Financial Precision: Uses integer math for all currency calculations to avoid floating-point inaccuracies, a standard practice in financial systems.
* Fully Tested: Includes a comprehensive suite of table-driven unit tests to ensure correctness and validate edge cases.

## Usage Example
The engine is designed to be simple to use. Here is how a service within the ERP would call the engine to process an inclusive amount:

```Go

package main

func main() {
    // 1. Define a complex tax template
    complexTemplate := model.TaxApplicationTemplate{
        Applications: []model.Application{
            {AppType: model.TaxTypePercent, Percent: 10, ApplyOn: model.TaxAppOnRefVal},
            {AppType: model.TaxTypeFixed, FixedAmount: 500},
            {AppType: model.TaxTypePercent, Percent: 5, ApplyOn: model.TaxAppOnRunningGross},
        },
    }

    // 2. Define an inclusive amount received from an invoice
    inclusiveAmount := 12075

    // 3. Calculate the exclusive base and tax
    exclusive, tax, err := model.CalculateWithInclusiveBase(inclusiveAmount, &complexTemplate)
    if err != nil {
        log.Fatal(err)
    }

    // Expected output: Exclusive Base: 10000, Tax: 2075
    fmt.Printf("Exclusive Base: %d, Tax: %d\n", exclusive, tax)
}
```

## Running the Demo
To run the main.go file which demonstrates a sample calculation:
```Bash
go run .
```

## Running the Tests
To run the full suite of unit tests for the model package:
```Bash
go test .
```
