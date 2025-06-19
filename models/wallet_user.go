package models

import "time"

type WalletUser struct {
	WalletID  int       `json:"wallet_id"`
	UserID    int       `json:"user_id"`
	UserRole  string    `json:"user_role"`
	CreatedAt time.Time `json:"created_at"`
}

type NewWalletUserRequest struct {
	WalletID int    `json:"wallet_id"`
	Username string `json:"username"`
	UserRole string `json:"user_role"`
}
