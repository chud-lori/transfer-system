package entities

import "github.com/shopspring/decimal"

type Transaction struct {
	Id                   int64
	SourceAccountID      int64
	DestinationAccountID int64
	Amount               decimal.Decimal
}
