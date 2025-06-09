package entities

type Account struct {
	AccountId int64  `json:"id"`
	Balance   string `json:"balance"`
}
