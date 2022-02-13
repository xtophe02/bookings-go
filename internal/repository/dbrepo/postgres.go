package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/xtophe02/bookings-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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
func (m *postgresDBRepo) GetUserByID(userID int) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, access_level ,created_at, updated_at FROM users WHERE id = $1;`

	row := m.DB.QueryRowContext(ctx, query, userID)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (m *postgresDBRepo) UpdateUserByID(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5;`

	_, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now())
	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, hashedPassword, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, hashedPassword, errors.New("wrong crendetials")
	} else if err != nil {
		return id, hashedPassword, err
	}
	return id, hashedPassword, nil
}
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, first_name, last_name, phone, email, start_date, end_date, room_id, r.created_at, r.updated_at, processed,
	rooms.id, rooms.room_name from reservations as r
	left join rooms on r.room_id = rooms.id
	order by start_date asc;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, nil
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Phone, &i.Email, &i.StartDate,
			&i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Processed, &i.Room.ID, &i.Room.RoomName)
		if err != nil {
			return reservations, nil
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select id, room_name, created_at, updated_at
	from rooms
	order by room_name;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, nil
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Room
		err := rows.Scan(&i.ID, &i.RoomName, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return rooms, nil
		}
		rooms = append(rooms, i)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, first_name, last_name, phone, email, start_date, end_date, room_id, r.created_at, r.updated_at ,
	rooms.id, rooms.room_name from reservations as r
	left join rooms on r.room_id = rooms.id
	where processed = 0
	order by start_date asc;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, nil
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Phone, &i.Email, &i.StartDate,
			&i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Room.ID, &i.Room.RoomName)
		if err != nil {
			return reservations, nil
		}
		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}
func (m *postgresDBRepo) GetReservationByID(reservationID int) (models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, first_name, last_name, phone, email, start_date, end_date, room_id, r.created_at, r.updated_at ,
	rooms.id, rooms.room_name, processed from reservations as r
	left join rooms on r.room_id = rooms.id
	where r.id = $1;`

	row := m.DB.QueryRowContext(ctx, query, reservationID)
	var i models.Reservation
	err := row.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Phone, &i.Email, &i.StartDate,
		&i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Room.ID, &i.Room.RoomName, &i.Processed)
	if err != nil {
		return i, err
	}

	return i, nil

}

func (m *postgresDBRepo) UpdateReservationByID(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5 where id=$6;`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) DeleteReservationByID(reservationID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id=$1;`

	_, err := m.DB.ExecContext(ctx, query, reservationID)
	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) UpdateProcessedForReservation(reservationID, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2;`

	_, err := m.DB.ExecContext(ctx, query, processed, reservationID)
	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `select id, coalesce(reservation_id,0), restriction_id, room_id, start_date, end_date
	from room_restrictions
	where $1 < end_date and $2 >= start_date and room_id = $3;`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, nil
	}
	defer rows.Close()
	for rows.Next() {
		var i models.RoomRestriction
		err := rows.Scan(&i.ID, &i.ReservationID, &i.RestrictionID, &i.RoomID, &i.StartDate, &i.EndDate)
		if err != nil {
			return restrictions, nil
		}
		restrictions = append(restrictions, i)
	}
	if err = rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil

}

func (m *postgresDBRepo) InsertBlockForRoom(roomID int, date time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id,
		created_at, updated_at) values($1,$2,$3,$4,$5,$6); `

	_, err := m.DB.ExecContext(ctx, query, date, date.AddDate(0, 0, 1), roomID, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (m *postgresDBRepo) DeleteBlockByID(roomRestrictionID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1; `

	_, err := m.DB.ExecContext(ctx, query, roomRestrictionID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
