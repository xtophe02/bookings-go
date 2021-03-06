package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"


	"github.com/xtophe02/bookings-go/internal/config"
)

var app *config.AppConfig

//TO HAVE IN THIS PKG THE APP.CONFIG AVAILABLE
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(rw http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(rw, http.StatusText(status), status)
}
func ServerError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool{
	return app.Session.Exists(r.Context(),"user_id")
}

