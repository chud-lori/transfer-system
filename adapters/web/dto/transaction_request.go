package dto

// @Description Transaction creation payload
type TransactionRequest struct {
	// @example 123
	SourceAccountID int64 `json:"source_account_id"`
	// @example 456
	DestinationAccountID int64 `json:"destination_account_id"`
	// @example 100.23344
	Amount string `json:"amount"`
}
