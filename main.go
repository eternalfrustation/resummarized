package main

import (
	"log"
	"net/http"


	"crypto/rand"

	"github.com/UniquityVentures/resummarized/pages"
	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	_ "github.com/joho/godotenv/autoload"
)



func main() {
	initFirebase()
	app := initApp()
	csrfKey := make([]byte, 128)
	rand.Read(csrfKey)
	csrfMiddleware := csrf.Protect(csrfKey, csrf.TrustedOrigins([]string{"localhost:7331", "localhost:4269"}), csrf.FieldName("_csrf"))

	router := getRoutes()
	log.Println("Listening on :4269")
	http.ListenAndServe(":4269", app.Handle(csrfMiddleware(router)))
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.Dashboard()).ServeHTTP(w, r)
}
