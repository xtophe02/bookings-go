package handlers

import (
	"encoding/gob"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
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

var functions = template.FuncMap{
	"humanDate":         render.HumanDate,
	"formatDate":        render.FormatDate,
	"iterate":           render.Iterate,
	"formatDateWeekDay": render.FormatDateWeekDay,
}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})
	gob.Register(models.Price{})
	gob.Register(map[string]int{})

	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//WE WILL RENDER ONCE ALL TEMPLATES
	tc, err := CreateTestTemplateCache()

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	if err != nil {
		log.Println("cannot create template cache")

	}
	app.TemplateCache = tc
	app.InProduction = true

	//CREATES REPOSITORY VARIABLE
	repo := NewTestingRepo(&app)
	//AFTER CREATING REPOSITORY, WE NEED TO PASS IT BACK TO HANDLER PKG
	NewHandlers(repo)
	//WE WILL GIVE RENDER PKG ACCESS TO THE MEMORY ADDRESS OF APPCONFIG
	render.NewRenderer(&app)

	os.Exit(m.Run())

}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/availability", Repo.Availability)
	mux.Post("/availability", Repo.PostAvailability)
	mux.Get("/reservation", Repo.Reservation)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/rooms/general-quarters", Repo.GeneralQuarters)
	mux.Get("/rooms/major-suite", Repo.MajorSuite)
	mux.Get("/user/login", Repo.Login)
	mux.Get("/user/logout", Repo.Logout)
	mux.Post("/user/login", Repo.PostLogin)
	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	mux.Post("/reservations-calendar", Repo.AdminPostReservationsCalendar)
	mux.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}", Repo.AdminPostShowReservation)
	mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation)
	mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pages []string

func walk(s string, d fs.DirEntry, err error) error {

	if err != nil {
		return err
	}
	if !d.IsDir() {

		pages = append(pages, s)
		// pages = append(pages,strings.TrimPrefix(s,"../../"))

	}
	return nil
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//MAP OF STRINGS AND VALUE OF TEXT-TEMPLATES
	myCache := map[string]*template.Template{}

	filepath.WalkDir("./../../templates", walk)

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
			return myCache, err
		}

		//CHECKS IF A LAYOUT EXISTS
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			//INJECTS THE LAYOUT ON TEXT-TEMPLATE
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	//RETURNS MAP OF NAMED FILED - TEXT-TEMPLATE

	return myCache, nil
}
