package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	IsDeposit bool      `json:"is_deposit"`
	Category  string    `json:"category"`
	WalletID  int       `json:"wallet_id"`
	CreatorID int       `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionInput struct {
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	WalletID int     `json:"wallet_id"`
}
