package repository

import (
	"ReservationsService/internal/core"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type DataBase struct {
	conn *pgx.Conn
}

func InitDataBase(ctx context.Context, connString string) (*DataBase, error) {
	dataBase := &DataBase{}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return dataBase, err
	}

	dataBase.conn = conn

	return dataBase, nil
}

func (dataBase *DataBase) Close(ctx context.Context) error {
	return dataBase.conn.Close(ctx)
}

func (dataBase *DataBase) CreateTable(ctx context.Context) error {
	strQuery := `
	CREATE TABLE IF NOT EXISTS reservations (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		author TEXT NOT NULL,
		title TEXT NOT NULL,
		taken_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		return_at TIMESTAMPTZ
	);
`

	_, err := dataBase.conn.Exec(ctx, strQuery)
	if err != nil {
		return err
	}

	return nil
}

func (dataBase *DataBase) AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error) {
	response := &core.ReservationResponse{ReservationStatus: false}

	strQuery := `
	INSERT INTO reservations (user_id, author, title, taken_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := dataBase.conn.Exec(ctx,
		strQuery,
		request.UserID,
		request.Author,
		request.Title,
		time.Now(),
	)

	if err != nil {
		return response, err
	}

	response.ReservationStatus = true

	return response, nil
}

func (dataBase *DataBase) CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error) {
	status := &core.ReturnResponse{Status: false}

	strQuery := `
	UPDATE reservations
		 SET return_at = NOW()
		 WHERE id = $1
`
	cmdTag, err := dataBase.conn.Exec(
		ctx,
		strQuery,
		request.ReservationID,
	)
	if err != nil {
		return status, err
	}

	if cmdTag.RowsAffected() == 0 {
		return status, fmt.Errorf("reservation not found or already returned")
	}

	status.Status = true

	return status, nil
}

func (dataBase *DataBase) CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error) {
	reservations := &core.CheckReservResponse{Reservations: []core.ReservationInfo{}}

	rows, err := dataBase.conn.Query(
		ctx,
		`SELECT id, user_id, author, title, taken_at, return_at
		 FROM reservations
		 WHERE user_id = $1
		 ORDER BY taken_at DESC`,
		userId,
	)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var info core.ReservationInfo

		err := rows.Scan(
			&info.ReservationID,
			&info.UserID,
			&info.Author,
			&info.Title,
			&info.TakenAt,
			&info.ReturnAt,
		)
		if err != nil {
			return reservations, err
		}

		reservations.Reservations = append(reservations.Reservations, info)
	}

	if rows.Err() != nil {
		return reservations, rows.Err()
	}

	return reservations, nil
}
