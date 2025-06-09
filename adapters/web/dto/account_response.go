package dto

type AccountResponse struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}
