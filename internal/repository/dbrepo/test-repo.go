package dbrepo

import (
	"errors"
	"time"

	"github.com/xtophe02/bookings-go/internal/models"
)

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}
func (m *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	if roomID > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
func (m *testDBRepo) GetUserByID(userID int) (models.User, error) {
	var user models.User

	return user, nil
}
func (m *testDBRepo) UpdateUserByID(user models.User) error {

	return nil
}
func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {
	if email == "me@here.ca" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")

}
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}
func (m *testDBRepo) GetReservationByID(reservationID int) (models.Reservation, error) {

	var reservation models.Reservation

	return reservation, nil
}
func (m *testDBRepo) UpdateReservationByID(u models.Reservation) error {

	return nil
}
func (m *testDBRepo) DeleteReservationByID(reservationID int) error {

	return nil
}
func (m *testDBRepo) UpdateProcessedForReservation(reservationID, processed int) error {

	return nil
}
func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil
}
func (m *testDBRepo) InsertBlockForRoom(roomID int, date time.Time) error {

	return nil
}
func (m *testDBRepo) DeleteBlockByID(roomRestrictionID int) error {

	return nil
}
