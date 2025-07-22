package main

import (
	"fmt"
	"log"
)

func main() {
	// Simple Template
	simpleTemplate := TaxApplicationTemplate{
		TemplateID: "Simple",
		Applications: []Application{
			{
				AppType: TaxTypePercent,
				Percent: 10.0,
				ApplyOn: TaxAppOnRefVal,
			},
		},
	}

	complexTemplate := TaxApplicationTemplate{
		TemplateID: "Complex",
		Applications: []Application{
			{
				AppType: TaxTypePercent,
				Percent: 10.0,
				ApplyOn: TaxAppOnRefVal,
			},
			{
				AppType:     TaxTypeFixed,
				FixedAmount: 2600, // In Cents
			},
			{
				AppType: TaxTypePercent,
				Percent: 10,
				ApplyOn: TaxAppOnRunningGross,
			},
		},
	}

	newInvoice := Invoice{
		Items: []LineItem{
			{
				ItemName:  "",
				Amount:    10000, // In Cents
				Inclusive: false,
				Template:  simpleTemplate,
			},
			{
				ItemName:  "",
				Amount:    22000, // In Cents
				Inclusive: true,
				Template:  complexTemplate,
			},
		},
	}

	cost, tax, entries, err := newInvoice.Process()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\nTotal Inventory Cost: %s\n", PrintMoney(cost))
	fmt.Printf("Total Tax Amount: %s\n", PrintMoney(tax))

	// Trial Balance
	totalDebits := 0
	totalCredits := 0

	fmt.Printf("\nLedger Entries:\n")
	for _, entry := range entries {
		totalDebits += entry.Debit
		totalCredits += entry.Credit
		fmt.Printf("Account: %s, \tDebit: %s, \tCredit: %s\n", entry.Account, PrintMoney(entry.Debit), PrintMoney(entry.Credit))
	}

	fmt.Printf("\n\tTotal Debits:\t%s\n", PrintMoney(totalDebits))
	fmt.Printf("\tTotal Credits:\t%s\n", PrintMoney(totalCredits))

}
