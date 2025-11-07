package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type App struct {
	Db *pgx.Conn
}

func DbConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func initApp() App {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, DbConnString())

	if err != nil {
		log.Printf("Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	init_posts, err := os.ReadFile("sql/init_posts.sql")
	if err != nil {
		log.Printf("Unable to find sql for creating posts table: %v\n", err)
		os.Exit(1)
	}
	conn.Exec(ctx, string(init_posts))
	return App{
		Db: conn,
	}
}

func (app *App) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "app", app)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
