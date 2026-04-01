package core

type CheckReservRequest struct {
	UserID int `json:"user_id"`
}

type CheckReservResponse struct {
	Reservations []ReservationInfo `json:"reservations"`
}
