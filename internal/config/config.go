package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

//APPCONFIG HOLDS THE APP CONFIG AND IT WILL BE AVAILABLE EVERYWHERE
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
}
