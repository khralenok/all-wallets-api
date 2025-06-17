package models

type WalletUsers struct {
	WalletID int    `json:"wallet_id"`
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}
