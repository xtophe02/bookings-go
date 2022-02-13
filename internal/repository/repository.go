package repository

import (
	"time"

	"github.com/xtophe02/bookings-go/internal/models"
)

//ALL FUNCS ON /dbrepo
type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(roomID int) (models.Room, error)
	GetUserByID(userID int) (models.User, error)
	UpdateUserByID(user models.User) error
	Authenticate(email, password string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(reservationID int) (models.Reservation, error)
	UpdateReservationByID(u models.Reservation) error
	DeleteReservationByID(reservationID int) error
	UpdateProcessedForReservation(reservationID, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	DeleteBlockByID(roomRestrictionID int) error
	InsertBlockForRoom(roomID int, date time.Time) error
}
