package handlers

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"github.com/xtophe02/bookings-go/internal/config"

	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager 
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{}

func getRoutes() http.Handler{
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	
	app.Session = session

	//WE WILL RENDER ONCE ALL TEMPLATES
	tc, err := CreateTestTemplateCache()
	
	if err != nil{
		log.Println("cannot create template cache")
	
	}
	app.TemplateCache = tc
	app.InProduction = true

	

	//CREATES REPOSITORY VARIABLE
	repo := NewRepo(&app)
	//AFTER CREATING REPOSITORY, WE NEED TO PASS IT BACK TO HANDLER PKG
	NewHandlers(repo)
		//WE WILL GIVE RENDER PKG ACCESS TO THE MEMORY ADDRESS OF APPCONFIG
	render.NewTemplates(&app)

	mux := chi.NewRouter()
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/",Repo.Home)
	mux.Get("/about",Repo.About)
	mux.Get("/availability",Repo.Availability)
	mux.Post("/availability",Repo.PostAvailability)
	mux.Get("/reservation",Repo.Reservation)
	mux.Post("/reservation",Repo.PostReservation)
	mux.Get("/contact",Repo.Contact)
	mux.Get("/reservation-summary",Repo.ReservationSummary)
	mux.Get("/rooms/general-quarters",Repo.GeneralQuarters)
	mux.Get("/rooms/major-suite",Repo.MajorSuite)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
	
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true, 
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}
//loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pages []string
func walk(s string, d fs.DirEntry, err error)  error {

	if err != nil {
		 return err
	}
	if ! d.IsDir() {
	
		pages = append(pages,s)
		// pages = append(pages,strings.TrimPrefix(s,"../../"))
		
	}
	return nil
}

func CreateTestTemplateCache() (map[string]*template.Template,error){
	//MAP OF STRINGS AND VALUE OF TEXT-TEMPLATES
	myCache :=map[string]*template.Template{}

	filepath.WalkDir("./../../templates",walk)

	// pages, err := filepath.Glob("./../../templates/*.page.tmpl")
	// if err != nil {
	// 	return myCache,err
	// }

	for _, page := range pages {
	
		//FETCH NAME OF FILE
		name := filepath.Base(page)
		//FETCH TEXT-TEMPLATE OF THE NAMED FILE
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		
		if err != nil {
			return myCache,err
		}

		//CHECKS IF A LAYOUT EXISTS
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
		if err != nil {
			return myCache,err
		}
	
		if len(matches) > 0{
			//INJECTS THE LAYOUT ON TEXT-TEMPLATE
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
			
			if err != nil {
				return myCache,err
			}
		}
	
		myCache[name] = ts
	}
	//RETURNS MAP OF NAMED FILED - TEXT-TEMPLATE 

	return myCache,nil
}