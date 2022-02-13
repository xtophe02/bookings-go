package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/driver"
	"github.com/xtophe02/bookings-go/internal/handlers"
	"github.com/xtophe02/bookings-go/internal/helpers"
	"github.com/xtophe02/bookings-go/internal/models"
	"github.com/xtophe02/bookings-go/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("start email listener...")
	listenForMail()

	fmt.Println("Starting app on port", portNumber)
	// http.ListenAndServe(portNumber,nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	//WHAT ARE WE GOING TO PUT IN THE SESSION
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})
	gob.Register(models.Price{})
	gob.Register(map[string]int{})

	//read flags
	inProduction := flag.Bool("production", true, "App is in production")
	useCache := flag.Bool("cache", true, "use template cache")
	dbName := flag.String("dbname", "", "Database Name")
	dbHost := flag.String("dbhost", "localhost", "Database Host")
	dbUser := flag.String("dbuser", "", "Database User")
	dbPass := flag.String("dbpass", "", "Database Password")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable,prefer,require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" || *dbPass == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	//CANNOT CLOSE IT HERE BECAUSE RUN RUNS ONCE
	app.MailChan = mailChan

	app.InProduction = *inProduction
	app.UseCache = *useCache

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//CONNECT TO DATABASE
	log.Println("connecting to database...")
	connectString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectString)
	if err != nil {
		log.Fatal("cannot connect to database!... leaving...")
	}
	log.Println("successful connection to database")
	//WE WILL RENDER ONCE ALL TEMPLATES
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	//WE WILL GIVE RENDER PKG ACCESS TO THE MEMORY ADDRESS OF APPCONFIG
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	//CREATES REPOSITORY VARIABLE
	repo := handlers.NewRepo(&app, db)
	//AFTER CREATING REPOSITORY, WE NEED TO PASS IT BACK TO HANDLER PKG
	handlers.NewHandlers(repo)

	return db, nil
}
