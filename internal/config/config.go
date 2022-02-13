package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/xtophe02/bookings-go/internal/models"
)

//APPCONFIG HOLDS THE APP CONFIG AND IT WILL BE AVAILABLE EVERYWHERE
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	MailChan      chan models.MailData
}
