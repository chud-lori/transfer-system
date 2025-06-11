package dto

import "github.com/shopspring/decimal"

// @Description Account creation payload
type AccountRequest struct {
	// Account ID
	// @example 123
	AccountID int64 `json:"account_id"`
	// Initial balance (string to allow decimal format)
	// @example 100.23344
	Balance string `json:"initial_balance"`
}

type InternalAccountRequest struct {
	AccountID int64
	Balance   decimal.Decimal
}
