package entities

import "github.com/shopspring/decimal"

type Account struct {
	AccountID int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
}
