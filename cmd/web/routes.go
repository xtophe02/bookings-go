package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/xtophe02/bookings-go/pkg/config"
	"github.com/xtophe02/bookings-go/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler{
	mux := chi.NewRouter()

	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about",http.HandlerFunc(handlers.Repo.About))

	return mux
}