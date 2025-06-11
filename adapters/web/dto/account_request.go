package dto

import "github.com/shopspring/decimal"

type AccountRequest struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"initial_balance"`
}

type InternalAccountRequest struct {
	AccountID int64
	Balance   decimal.Decimal
}
