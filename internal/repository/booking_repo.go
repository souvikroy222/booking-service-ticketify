package repository

import (
	"booking-ticketify/internal/models"
	"context"
	"database/sql"
)

type BookingRepository struct {
	DB *sql.DB
}

// create the constructor here
func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, req models.BookingRequest) {
	// query :=`
	// 		INSERT INTO bookings (user_id,seat_id,status)
	// 		VALUES ($1,$2, '')
	// `
}
