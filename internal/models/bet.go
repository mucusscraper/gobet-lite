package models

import "time"

type Bet struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Input      float64   `json:"input"`
	Multiplier float64   `json:"multiplier"`
	Result     float64   `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateBetRequest struct {
	UserID     int     `json:"user_id"`
	Input      float64 `json:"input"`
	Multiplier float64 `json:"multiplier"`
}
