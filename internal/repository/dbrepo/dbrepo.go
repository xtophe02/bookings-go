package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//CREATE A DATABASE REPOSITORY SO THE APP HAS ACCCESS ANYWHERE
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{App: a, DB: conn}
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//IF SOMETHING WRONG WITH THE USER... CANCEL IN 3S THE QUERY
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	query := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
						VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err := m.DB.QueryRowContext(ctx, query, res.FirstName, res.LastName, res.Email,
		res.Phone, res.StartDate, res.EndDate,
		res.RoomID, time.Now(), time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id ,created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7);`

	_, err := m.DB.ExecContext(ctx, query, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, res.RestrictionID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	var numRows int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count (id) from room_restrictions where 
	room_id = $1 and
	$2 < end_date and $3 > start_date;`
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)

	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.room_name 
	from rooms as r 
	where r.id 
	not in 
		(select rr.room_id 
		 from room_restrictions as rr 
		 where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}
func (m *postgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, room_name, created_at, updated_at FROM rooms WHERE id = $1;`

	row := m.DB.QueryRowContext(ctx, query, roomID)

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
