package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/handlers"
	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/render"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main(){
	err := run()
if err != nil{
	log.Fatal(err)
}
	fmt.Println("Starting app on port",portNumber)
	// http.ListenAndServe(portNumber,nil)
	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}
}

func run() error {
		//WHAT ARE WE GOING TO PUT IN THE SESSION
		gob.Register(models.Reservation{})

		app.InProduction = false
	
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
		
		app.Session = session
	
		//WE WILL RENDER ONCE ALL TEMPLATES
		tc, err := render.CreateTemplateCache()
		if err != nil{
			log.Println("cannot create template cache")
			return err
		}
		app.TemplateCache = tc
		app.UseCache = false
		//WE WILL GIVE RENDER PKG ACCESS TO THE MEMORY ADDRESS OF APPCONFIG
		render.NewTemplates(&app)
	
		//CREATES REPOSITORY VARIABLE
		repo := handlers.NewRepo(&app)
		//AFTER CREATING REPOSITORY, WE NEED TO PASS IT BACK TO HANDLER PKG
		handlers.NewHandlers(repo)
	return nil
}