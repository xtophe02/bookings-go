package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler{
	mux := chi.NewRouter()
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
	mux.Get("/availability",http.HandlerFunc(handlers.Repo.Availability))
	mux.Post("/availability",http.HandlerFunc(handlers.Repo.PostAvailability))
	mux.Get("/reservation",http.HandlerFunc(handlers.Repo.Reservation))
	mux.Post("/reservation",http.HandlerFunc(handlers.Repo.PostReservation))
	mux.Get("/contact",http.HandlerFunc(handlers.Repo.Contact))
	mux.Get("/reservation-summary",http.HandlerFunc(handlers.Repo.ReservationSummary))
	mux.Get("/rooms/general-quarters",http.HandlerFunc(handlers.Repo.GeneralQuarters))
	mux.Get("/rooms/major-suite",http.HandlerFunc(handlers.Repo.MajorSuite))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))
	return mux
}