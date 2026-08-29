package models

import "time"

type BookingRequest struct {
	UserID int    `json:"user_id"`
	SeatId string `json:"seat_id"`
}
type Booking struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	SeatId    string    `json:"seat_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
