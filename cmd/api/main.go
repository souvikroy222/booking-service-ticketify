package main

import (
	"booking-ticketify/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {

	//add the redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//check the redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect with redis")
	}
	log.Println("Connected to redis")

	r := chi.NewRouter()

	//use built in middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/seats/select", func(w http.ResponseWriter, r *http.Request) {
		var bookingReq models.BookingRequest
		err := json.NewDecoder(r.Body).Decode(&bookingReq)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		//marshal the data

		ctx := r.Context()
		newSeatKey := fmt.Sprintf("soft_lock:%s", bookingReq.SeatId)
		userLockKey := fmt.Sprintf("user_lock:%s", bookingReq.UserID)

		oldSeatId, err := rdb.Get(ctx, userLockKey).Result()
		if err == nil && oldSeatId != "" {

			if oldSeatId == bookingReq.SeatId {
				fmt.Println("Same seat trying to book")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("you are already booked this seat"))
				return
			}

			oldSeatKey := fmt.Sprintf("soft_lock:%s", oldSeatId)
			log.Printf("Released the occupied seat of %s", oldSeatKey)
			rdb.Del(ctx, oldSeatKey)

		}

		//try to set for 2 minutes
		success, err := rdb.SetNX(ctx, newSeatKey, bookingReq.UserID, 120*time.Second).Result()
		if err != nil {
			log.Printf("Failed to execute SET NX %v", err)
		}

		if !success {
			fmt.Println("Could not lock the seat ID, because it already exist")
		}
		rdb.Set(ctx, userLockKey, bookingReq.SeatId, 120*time.Second)

		fmt.Printf("Received data %s ", bookingReq.SeatId)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("JSON received successfully"))
	})

	log.Println("Ticketing server is running")

	http.ListenAndServe(":8080", r)

}
