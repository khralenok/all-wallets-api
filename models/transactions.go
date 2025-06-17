package models

import "time"

type Transactions struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	IsDeposit bool      `json:"is_deposit"`
	Category  string    `json:"category"`
	WalletID  int       `json:"wallet_id"`
	CreatorID int       `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}
