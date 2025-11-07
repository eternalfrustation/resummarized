package main

import (
	"net/http"

	"github.com/UniquityVentures/resummarized/middlewares"
	"github.com/UniquityVentures/resummarized/pages"
	"github.com/a-h/templ"
	"github.com/rickb777/servefiles/v3"
)

func getRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/", templ.Handler(pages.HomePage()))
	router.Handle("/assets/", servefiles.NewAssetHandler("./assets/").StripOff(1) /*.WithMaxAge(24 * time.Hour) */)
	router.Handle("/dashboard", middlewares.FirebaseAuthMiddleware(authClient, http.HandlerFunc(DashboardHandler)))
	router.Handle("/auth/sessionLogin/", middlewares.FormHandler(SessionLoginHandler))
	return router
}
