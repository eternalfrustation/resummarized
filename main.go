package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"crypto/rand"

	"github.com/gorilla/csrf"
	_ "github.com/joho/godotenv/autoload"
)

var app App

func init() {
	initFirebase()
	app = initApp()
	flag.Func("add_admin", "used to specify email of admin to be added", func(value string) error {
		return app.ExecQuery(context.Background(), "sql/create_admin.sql", value)
	})
	flag.Func("remove_admin", "used to specify email of admin to be removed", func(value string) error {
		return app.ExecQuery(context.Background(), "sql/remove_admin.sql", value)
	})
	flag.BoolFunc("list_admin", "used to fetch a list of admins", func(value string) error {
		admins, err := FetchRows[Admin](&app, context.Background(), "sql/list_admin.sql")
		if err != nil {
			return err
		}
		fmt.Printf("Admin Id\t-\tAdmin Email\n")
		fmt.Printf(strings.Repeat("-", 2 * len("Admin Id\tAdmin Email")))
		fmt.Print("\n")
		for _, admin := range admins {
			fmt.Printf("%d\t\t-\t%s\n", admin.AdminID, admin.AdminEmail)
		}
		return nil
	})
}

func main() {
	flag.Parse()
	csrfKey := make([]byte, 128)
	rand.Read(csrfKey)
	csrfMiddleware := csrf.Protect(csrfKey, csrf.TrustedOrigins([]string{"localhost:7331", "localhost:426"}), csrf.FieldName("_csrf"))

	router := getRoutes()
	log.Println("Listening on :426")
	http.ListenAndServe(":4269", app.Handle(csrfMiddleware(router)))
}
