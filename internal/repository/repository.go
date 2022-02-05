package repository

import (
	"time"

	"github.com/xtophe02/bookings-go/internal/models"
)

//ALL FUNCS ON /dbrepo
type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(roomID int) (models.Room, error)
}
