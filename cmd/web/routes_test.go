package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/xtophe02/bookings-go/internal/config"
)

func TestRoutes(t *testing.T){
	var app config.AppConfig

	mux:=routes(&app)

	switch v := mux.(type){
	case *chi.Mux:
		//do nothing
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, is %T",v))
	}
}