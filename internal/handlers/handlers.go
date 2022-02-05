package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/driver"
	"github.com/xtophe02/bookings-go/internal/forms"
	"github.com/xtophe02/bookings-go/internal/helpers"
	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/render"
	"github.com/xtophe02/bookings-go/internal/repository"
	"github.com/xtophe02/bookings-go/internal/repository/dbrepo"
)

// REPOSITORY TYPE
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//REPO REPOSITORY USED BY THE HANDLERS
var Repo *Repository

//CREATES NEW REPOSITORY
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: a, DB: dbrepo.NewPostgresRepo(db.SQL, a)}
}

//SETS THE REPOSITORY FOR THE HANDLERS
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	// remoteIP := r.RemoteAddr
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(rw, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello from handlers"
	// remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = remoteIP
	render.Template(rw, r, "about.page.tmpl", &models.TemplateData{
		// StringMap: stringMap,
	})
}
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(rw, r, "contact.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	//NEED OF CSRFTOKEN.. INJECTED ON Template AS TEMPLATEDATE DEFAULT

	render.Template(rw, r, "availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {

		helpers.ServerError(rw, err)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	//YYYY-MM-DD -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		//TODO: send flash error if no dates
		helpers.ServerError(rw, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No avaulability")
		http.Redirect(rw, r, "/availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	render.Template(rw, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	// rw.Write([]byte(fmt.Sprintf("start data is %v and end is %v", sd, ed)))

}
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {

	//NEED TO CREATE A DATA-RESERVATION TO REPOPULATE IF NEEDED
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("couldn't get the session"))
	}
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	reservation.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", reservation)
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(rw, r, "reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("counld get session reservation"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	// sd := r.Form.Get("start_date")
	// ed := r.Form.Get("end_date")

	//YYYY-MM-DD -- 01/02 03:04:05PM '06 -0700

	// layout := "2006-01-02"
	// startDate, err := time.Parse(layout, sd)
	// if err != nil {
	// 	helpers.ServerError(rw, err)
	// 	return
	// }
	// endDate, err := time.Parse(layout, ed)
	// if err != nil {
	// 	helpers.ServerError(rw, err)
	// 	return
	// }
	// roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	// if err != nil {
	// 	helpers.ServerError(rw, err)
	// }
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	// reservation := models.Reservation{
	// 	FirstName: r.Form.Get("first_name"),
	// 	LastName:  r.Form.Get("last_name"),
	// 	Email:     r.Form.Get("email"),
	// 	Phone:     r.Form.Get("phone"),
	// 	StartDate: startDate,
	// 	EndDate:   endDate,
	// 	RoomID:    roomID,
	// }
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")

	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(rw, r, "reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//SAVE TO DB
	newReservationdID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	// m.App.Session.Put(r.Context(),"reservation",reservation)
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationdID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	//WE NEED TO TELL THE SESSION WHAT TYPE ARE WE STORING
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}
func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("couldn't get session"))
	}
	reservation.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation", http.StatusSeeOther)

}
func (m *Repository) BookRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	sd := r.URL.Query().Get("s")
	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		helpers.ServerError(rw, errors.New("could not get end date"))
		return
	}
	ed := r.URL.Query().Get("e")
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		helpers.ServerError(rw, errors.New("could not get end date"))
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	var reservation models.Reservation
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.RoomID = roomID
	reservation.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation", http.StatusSeeOther)

}
func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(rw, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
func (m *Repository) GeneralQuarters(rw http.ResponseWriter, r *http.Request) {
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello from handlers"
	// remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = remoteIP
	render.Template(rw, r, "general-quarters.page.tmpl", &models.TemplateData{
		// StringMap: stringMap,
	})
}
func (m *Repository) MajorSuite(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(rw, r, "major-suite.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start_date")
	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		helpers.ServerError(rw, errors.New("could not get start date"))
		return
	}
	ed := r.Form.Get("end_date")
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		helpers.ServerError(rw, errors.New("could not get end date"))
		return
	}
	room_id := r.Form.Get("room_id")
	roomId, err := strconv.Atoi(room_id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	resp := jsonResponse{available, "", room_id, sd, ed}
	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}
