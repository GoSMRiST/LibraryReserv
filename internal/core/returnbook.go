package core

type ReturnRequest struct {
	ReservationID int `json:"reservation_id"`
}

type ReturnResponse struct {
	Status bool `json:"status"`
}
