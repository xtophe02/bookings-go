package handlers

import (
	"net/http"

	"github.com/xtophe02/bookings-go/pkg/config"
	"github.com/xtophe02/bookings-go/pkg/models"
	"github.com/xtophe02/bookings-go/pkg/render"
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
	render.RenderTemplate(rw, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository)About(rw http.ResponseWriter, r *http.Request){
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from handlers"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(rw, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
