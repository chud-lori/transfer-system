package entities

type Transaction struct {
	Id                   int64  `json:"id"`
	SourceAccountId      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}
