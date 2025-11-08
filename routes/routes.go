package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/UniquityVentures/resummarized/handlers"
	"github.com/UniquityVentures/resummarized/middlewares"
	"github.com/UniquityVentures/resummarized/pages"
	"github.com/a-h/templ"
	"github.com/rickb777/servefiles/v3"
)

func GetRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/assets/", servefiles.NewAssetHandler("./assets/").StripOff(1).WithMaxAge(24*time.Hour))
	router.Handle("/auth/sessionLogin/", middlewares.FormHandler(handlers.SessionLoginHandler))
	router.Handle("/", templ.Handler(pages.HomePage()))
	Nest(router, "/user", getUserRoutes())
	Nest(router, "/admin", getAdminRoutes())
	return router
}

func Nest(router *http.ServeMux, route string, h http.Handler) {
	router.Handle(fmt.Sprintf("%s/", route), http.StripPrefix(route, h))

}

func getAdminRoutes() *http.ServeMux {
	router := http.NewServeMux()
	return router
}

func getUserRoutes() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
