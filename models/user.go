package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"user_pwd"`
	BaseCurrency string    `json:"base_currency"`
	CreatedAt    time.Time `json:"created_at"`
	IsDeleted    bool      `json:"is_deleted"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type UserOutput struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	BaseCurrency string `json:"base_currency"`
}

type LoginInputs struct {
	Username string `json:"username"`
	Password string `json:"user_pwd"`
}

type SigninInputs struct {
	Username     string `json:"username"`
	Password     string `json:"user_pwd"`
	BaseCurrency string `json:"base_currency"`
}
