package models

import "time"

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
	// UserID    int
	Room      Room
	CreatedAt time.Time
	UpdatedAt time.Time
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
}

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type RoomRestriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartDate       time.Time
	EndDate         time.Time
	RoomID          int
	Room            Room
	ReservationID   int
	Reservation     Reservation
	RestrictionID   int
	Restriction     Restriction
}

type Price struct {
	ID          int
	RoomID      int
	Room        Room
	WinterPrice int
	SummerPrice int
}
