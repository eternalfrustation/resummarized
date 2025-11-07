package main

import (
	"net/http"

	"github.com/UniquityVentures/resummarized/pages"
	"github.com/a-h/templ"
)


func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.Dashboard()).ServeHTTP(w, r)
}
