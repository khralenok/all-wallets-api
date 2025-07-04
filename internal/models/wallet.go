package models

import "time"

type Wallet struct {
	ID           int       `json:"id"`
	WalletName   string    `json:"wallet_name"`
	Currency     string    `json:"currency"`
	Balance      int       `json:"balance"`
	LastSnapshot time.Time `json:"last_snapshot"`
	CreatedAt    time.Time `json:"created_at"`
}

type NewWalletRequest struct {
	WalletName string `json:"wallet_name"`
	Currency   string `json:"currency"`
}

type WalletOutput struct {
	WalletID            int    `json:"wallet_id"`
	WalletName          string `json:"wallet_name"`
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	UserCurrencyBalance string `json:"user_currency_balance"`
	UserRole            string `json:"user_role"`
}
