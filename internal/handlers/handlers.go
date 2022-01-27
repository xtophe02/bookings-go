package handlers

import (
	"log"
	"net/http"

	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/forms"
	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/render"
)

// REPOSITORY TYPE
type Repository struct {
	App *config.AppConfig
}
//REPO REPOSITORY USED BY THE HANDLERS
var Repo * Repository

//CREATES NEW REPOSITORY
func NewRepo (a *config.AppConfig) *Repository{
	return &Repository{ App: a}
}

//SETS THE REPOSITORY FOR THE HANDLERS
func NewHandlers(r *Repository){
	Repo = r
}

func (m *Repository)Home(rw http.ResponseWriter, r *http.Request){
remoteIP := r.RemoteAddr
m.App.Session.Put(r.Context(),"remote_ip",remoteIP)
	render.RenderTemplate(rw, r,"home.page.tmpl", &models.TemplateData{})
}
func (m *Repository)About(rw http.ResponseWriter, r *http.Request){
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(rw, r,"about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (m *Repository)Contact(rw http.ResponseWriter, r *http.Request){
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(rw, r,"contact.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (m *Repository)Availability(rw http.ResponseWriter, r *http.Request){
	//NEED OF CSRFTOKEN.. INJECTED ON RENDERTEMPLATE AS TEMPLATEDATE DEFAULT

	render.RenderTemplate(rw, r,"availability.page.tmpl", &models.TemplateData{})
}


func (m *Repository)PostAvailability(rw http.ResponseWriter, r *http.Request){

}
func (m *Repository)Reservation(rw http.ResponseWriter, r *http.Request){

	//NEED TO CREATE A DATA-RESERVATION TO REPOPULATE IF NEEDED
	var emptyRervation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyRervation

	render.RenderTemplate(rw, r,"reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
func (m *Repository)PostReservation(rw http.ResponseWriter, r *http.Request){
	err:= r.ParseForm()
	if err != nil{
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
	}	
	form := forms.New(r.PostForm)

	form.Required("first_name","last_name","email")

	form.MinLength("first_name",3,r)
	form.IsEmail("email")

	if !form.Valid(){
		
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(rw, r,"reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//WE NEED TO TELL THE SESSION WHAT TYPE ARE WE STORING
	m.App.Session.Put(r.Context(), "reservation",reservation)
	http.Redirect(rw, r, "/reservation-summary",http.StatusSeeOther)
}
func (m *Repository)ReservationSummary(rw http.ResponseWriter, r *http.Request){
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(),"error","Can't get reservation from session")
		http.Redirect(rw,r,"/",http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(),"reservation")
	
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(rw, r,"reservation-summary.page.tmpl", &models.TemplateData{
	Data: data,
	})
}
func (m *Repository)GeneralQuarters(rw http.ResponseWriter, r *http.Request){
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(rw, r,"general-quarters.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (m *Repository)MajorSuite(rw http.ResponseWriter, r *http.Request){
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(rw, r,"major-suite.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
