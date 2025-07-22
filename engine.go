package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (inv *Invoice) Process() (invCost, tax int, entries []LedgerEntry, er error) {

	for _, li := range inv.Items {

		// Correct base Amount for Line Item
		if li.Inclusive {
			exBase, _, err := CalculateWithInclusiveBase(li.Amount, &li.Template)
			if err != nil {
				er = err
				return
			}

			li.Amount = exBase
		}

		// Build Results
		t, es, err := li.Process()
		if err != nil {
			er = err
			return
		}

		invCost += li.Amount
		tax += t
		entries = append(entries, es...)

	}

	return

}

func (li *LineItem) Process() (tax int, entries []LedgerEntry, er error) {

	tax, err := TaxCaculate(li.Amount, &li.Template)
	if err != nil {
		er = err
		return
	}

	// Build Entries
	entries = append(entries, []LedgerEntry{
		{
			Account: "Tax Asset",
			Debit:   tax,
		},
		{
			Account: "Supplier",
			Credit:  li.Amount + tax,
		},
		{
			Account: "Inventory",
			Debit:   li.Amount,
		},
	}...)

	return

}

// Finds the exclusive base, by iterative estimation algorithm, with possibly 1-cent rounding tolerance ...
func CalculateWithInclusiveBase(inclusiveBase int, template *TaxApplicationTemplate) (exclusiveBase int, tax int, er error) {

	if inclusiveBase <= 0 {
		// tax cannot be applied
		er = fmt.Errorf("validation error: Cannot calculate Tax for '0' or 'Negative' dollars as inclusive of tax amount")
		return
	}

	// estimateExclusive exlusive base
	estimateExclusive := 0
	taxEstimate := 0

	var totalRate float64 = 0
	var totalFixedApplied int = 0

	for _, app := range template.Applications {
		if app.AppType == TaxTypePercent {
			totalRate += app.Percent
		}
		if app.AppType == TaxTypeFixed {
			totalFixedApplied += app.FixedAmount
		}
	}

	// Estimat Base Amount
	estimateExclusive = inclusiveBase - Percent(inclusiveBase-totalFixedApplied, totalRate) - totalFixedApplied
	fmt.Printf("First Exclusive Estimate: %v\n", estimateExclusive)

	diff := 0
	for {
		// Run calc until result + diff == goal
		calcTax, err := TaxCaculate(estimateExclusive, template)
		if err != nil {
			er = err
			return
		}

		diff = (estimateExclusive + calcTax) - inclusiveBase
		fmt.Println("Diff: ", diff)
		taxEstimate = calcTax // Update Tax

		if diff >= -1 && diff <= 1 {
			// Try with 0 diff once
			finalGuess := estimateExclusive - diff
			finalTax, err := TaxCaculate(finalGuess, template)
			if err != nil {
				er = fmt.Errorf("final Tax Calculation failed: %v", err.Error())
				return
			}

			if (finalGuess + finalTax) == inclusiveBase {
				estimateExclusive = finalGuess
				taxEstimate = finalTax
				diff = 0
				break
			}
			break
		}
		estimateExclusive = estimateExclusive - diff
	}

	// Now Exclusive estimate will only have a difference of 1 cent
	// add difference to tax
	exclusiveBase = estimateExclusive
	tax = taxEstimate - diff

	return

}

func TaxCaculate(baseAmount int, template *TaxApplicationTemplate) (int, error) {

	if baseAmount < 0 {
		// Cannot apply Tax
		return 0, fmt.Errorf("validation error: Cannot calculate Tax for '0' or 'Negative' dollars as base amount")
	}

	// Apply On Logic
	runningGross := baseAmount
	prevTax := 0

	for _, app := range template.Applications {

		if app.AppType == "" {
			return 0, fmt.Errorf("application type not defined on template")
		}

		switch app.AppType {

		case TaxTypePercent:
			// Apply Percentage and Update Control vars
			switch app.ApplyOn {
			case TaxAppOnRefVal:
				tax := Percent(baseAmount, app.Percent)
				prevTax = tax
				runningGross += tax
				// Also Create the tax Application Record

			case TaxAppOnPrevTax:
				tax := Percent(prevTax, app.Percent)
				prevTax = tax
				runningGross += tax
				// Also Create the tax Application Record

			case TaxAppOnRunningGross:
				tax := Percent(runningGross, app.Percent)
				prevTax = tax
				runningGross += tax
			}

		case TaxTypeFixed:
			runningGross += app.FixedAmount
		default:
			return 0, fmt.Errorf("validation error: Invalid Application Type")
		}
	}

	return runningGross - baseAmount, nil

}

// Applies Round Up Method, common in Financial Systems
// Uses Integer math to avoid Floating Point arithmetic
func Percent(mnt int, rate float64) int {

	truncRate := math.Trunc(rate*100) / 100
	basePoints := int(math.Round(truncRate * 100))
	return (mnt*basePoints + 5000) / 10000

}

func PrintMoney(mnt int) string {

	str := strconv.Itoa(mnt)
	slc := strings.Split(str, "")
	switch {
	case len(slc) == 0:
		return "0.00"
	case len(slc) == 1:
		return "0.0" + strings.Join(slc, "")
	case len(slc) == 2:
		return "0." + strings.Join(slc, "")
	default:
		return strings.Join(slc[:len(slc)-2], "") + "." + strings.Join(slc[len(slc)-2:], "")
	}

}
