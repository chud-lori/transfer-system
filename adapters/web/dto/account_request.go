package dto

import "github.com/shopspring/decimal"

// type CreateAccountRequest struct {
// 	AccountID      int64  `json:"account_id"`
// 	InitialBalance string `json:"initial_balance"`
// }

type AccountRequest struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"initial_balance"`
}

type InternalAccountRequest struct {
	AccountID int64
	Balance   decimal.Decimal
}
