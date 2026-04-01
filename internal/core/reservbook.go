package core

import "time"

type ReservationRequest struct {
	UserID int    `json:"user_id"`
	Author string `json:"author"`
	Title  string `json:"title"`
}

type ReservationResponse struct {
	ReservationStatus bool `json:"reservation_status"`
}

type ReservationInfo struct {
	ReservationID int
	UserID        int
	Author        string
	Title         string
	TakenAt       time.Time
	ReturnAt      *time.Time
}
