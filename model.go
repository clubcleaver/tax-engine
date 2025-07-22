package main

// Taxation Model
type TaxCalculationType string

const (
	TaxTypePercent TaxCalculationType = "Percent"
	TaxTypeFixed   TaxCalculationType = "Fixed"
)

type TaxApplyOn string

const (
	TaxAppOnRefVal       TaxApplyOn = "Reference Value"
	TaxAppOnRunningGross TaxApplyOn = "Running Gross"
	TaxAppOnPrevTax      TaxApplyOn = "Previous Tax Amount"
)

type Application struct {
	AppType     TaxCalculationType
	FixedAmount int
	Percent     float64
	ApplyOn     TaxApplyOn
}

// TaxApplicationTemplate defines a set of tax rules to be applied sequentially to a value.
type TaxApplicationTemplate struct {
	TemplateID   string
	Applications []Application
}

// Ledger Record Example
type LedgerEntry struct {
	Account string
	Debit   int
	Credit  int
}

// Invoice
type LineItem struct {
	ItemName  string
	Amount    int
	Inclusive bool
	Template  TaxApplicationTemplate
}

type Invoice struct {
	Items []LineItem
}
