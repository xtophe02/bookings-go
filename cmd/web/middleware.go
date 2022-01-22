package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(rw, r)
// 	})
// }
//adds Cross Site Request Forgery to all POST requests
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
