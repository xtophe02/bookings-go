package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
func NewTestingRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a, DB: dbrepo.NewTestingRepo(a)}
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

	//good pratice to test the parse of a form
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
		// helpers.ServerError(rw, errors.New("couldn't get the session"))
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		// helpers.ServerError(rw, err)
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
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

	// good pratice
	err := r.ParseForm()
	if err != nil {
		// helpers.ServerError(rw, err)
		// return
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// helpers.ServerError(rw, errors.New("counld get session reservation"))
		m.App.Session.Put(r.Context(), "error", "can't get session reservation")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
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
		http.Error(rw, "Form not valid", http.StatusSeeOther)
		render.Template(rw, r, "reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//SAVE TO DB
	newReservationdID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		// helpers.ServerError(rw, err)
		m.App.Session.Put(r.Context(), "error", "can't insert data on db")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
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
		// helpers.ServerError(rw, err)
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction on db")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s, <br>
	This is confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "admin@christophemoreira.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	htmlMessage = fmt.Sprintf(`
	<strong>Reservation Notification</strong><br>
	A reservation has been made for %s from %s to %s
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.MailData{
		To:      "admin@christophemoreira.com",
		From:    "admin@christophemoreira.com",
		Subject: "Reservation Notification",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

	//WE NEED TO TELL THE SESSION WHAT TYPE ARE WE STORING
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}
func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	//hard to test
	// roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
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
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(out)
		return
	}

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
		resp := jsonResponse{
			Ok:      false,
			Message: "Error connecting to db",
		}
		out, _ := json.MarshalIndent(resp, "", "   ")
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(out)
		return
	}

	resp := jsonResponse{available, "", room_id, sd, ed}
	out, _ := json.MarshalIndent(resp, "", "   ")

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}
func (m *Repository) Login(rw http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["login"] = models.User{
		Email:    "",
		Password: "",
	}
	render.Template(rw, r, "login.page.tmpl", &models.TemplateData{Form: forms.New(nil), Data: data})
}
func (m *Repository) Logout(rw http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
func (m *Repository) PostLogin(rw http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())
	// good pratice
	err := r.ParseForm()
	if err != nil {
		// helpers.ServerError(rw, err)
		// return
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["login"] = models.User{
			Email:    "",
			Password: "",
		}
		// http.Error(rw, "Form not valid", http.StatusSeeOther)
		render.Template(rw, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	id, _, err := m.DB.Authenticate(r.Form.Get("email"), r.Form.Get("password"))
	if err != nil {

		data := make(map[string]interface{})
		data["login"] = models.User{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
		m.App.Session.Put(r.Context(), "error", "Invalid Login Credentials")
		render.Template(rw, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)

	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminNewReservations(rw http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(rw, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminAllReservations(rw http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	// fmt.Println(reservations)
	render.Template(rw, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminReservationsCalendar(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	LastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = LastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(rw, err)
			return
		}

		for _, y := range restrictions {

			if y.ReservationID > 0 {
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(rw, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
func (m *Repository) AdminPostShowReservation(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservationByID(reservation)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "changes saved")
	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", url[3]), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}
func (m *Repository) AdminPostReservationsCalendar(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	form := forms.New(r.PostForm)

	//remove blocks
	for _, x := range rooms {
		curBlockMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curBlockMap {
			if val, ok := curBlockMap[name]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}

	//add blocks
	for name := range r.PostForm {
		// fmt.Println(name)
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])

			// fmt.Println(roomID, exploded[3])
			t, _ := time.Parse("2006-01-2", exploded[3])
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", r.Form.Get("y"), r.Form.Get("m")), http.StatusSeeOther)
}
func (m *Repository) AdminShowReservation(rw http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	stringMap := make(map[string]string)
	stringMap["src"] = url[3]

	stringMap["year"] = r.URL.Query().Get("y")
	stringMap["month"] = r.URL.Query().Get("m")

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	// log.Print(reservation)
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(rw, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}
func (m *Repository) AdminProcessReservation(rw http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as Processed")

	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", url[3]), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}
func (m *Repository) AdminDeleteReservation(rw http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	err = m.DB.DeleteReservationByID(id)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", url[3]), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}
