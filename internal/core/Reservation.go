package core

import "time"

type ReservationRequest struct {
	UserID int
	Author string
	Title  string
}

type ReservationResponse struct {
	ReservationStatus bool
}

type ReturnRequest struct {
	UserID int
	Author string
	Title  string
}

type ReturnResponse struct {
	Status bool
}

type ReservationInfo struct {
	ReservationID int
	UserID        int
	Author        string
	Title         string
	TakenAt       time.Time
	ReturnAt      *time.Time
}

type CheckReservResponse struct {
	Reservations []ReservationInfo
}

type CheckAvailabilityRequest struct {
	Author string
	Title  string
}

type CheckAvailabilityResponse struct {
	Result bool
}
