package dto

import "time"

type Order struct {
	UserID      string
	OrderNumber string
	OrderDate   string
}

type FullOrder struct {
	OrderNumber string  `json:"number"`
	OrderStatus string  `json:"status"`
	Accrual     float32 `json:"accrual"`
	OrderDate   string  `json:"uploaded_at"`
}

type UserBalance struct {
	Current   float32
	Withdrawn float32
}

type WithdrawOrder struct {
	Order string
	Sum   float32
}

type Billing struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
