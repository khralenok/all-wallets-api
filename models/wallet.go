package models

import "time"

type Wallet struct {
	ID           int       `json:"id"`
	WalletName   string    `json:"wallet_name"`
	Currency     string    `json:"currency"`
	Balance      string    `json:"balance"`
	LastSnapshot time.Time `json:"last_snapshot"`
	CreatedAt    time.Time `json:"created_at"`
}
