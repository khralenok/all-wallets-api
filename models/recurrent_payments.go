package models

import "time"

type RecurrentPayment struct {
	ID               int       `json:"id"`
	Amount           int       `json:"amount"`
	IsDeposit        bool      `json:"is_deposit"`
	Category         string    `json:"category"`
	WalletID         int       `json:"wallet_id"`
	Frequency        string    `json:"frequency"`
	ScheduledDay     int       `json:"scheduled_day"`
	ScheduledWeekday int       `json:"scheduled_weekday"`
	ScheduledMonth   int       `json:"scheduled_month"`
	NextRun          time.Time `json:"next_run"`
	EndAt            time.Time `json:"end_at"`
	CreatorID        int       `json:"creator_id"`
	CreatedAt        time.Time `json:"created_at"`
}
