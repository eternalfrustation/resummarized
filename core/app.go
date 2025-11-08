package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Db *pgxpool.Pool
}

func DbConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func InitApp() App {
	ctx := context.Background()
	poolConfig, err := pgxpool.ParseConfig(DbConnString())
	if err != nil {
		log.Printf("Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		log.Printf("Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	app := App{
		Db: conn,
	}
	if err = app.ExecQuery(ctx, "sql/init_articles.sql"); err != nil {
		log.Fatal(err)
	}
	if err = app.ExecQuery(ctx, "sql/init_admin.sql"); err != nil {
		log.Fatal(err)
	}

	return app
}

func (app *App) ExecQuery(ctx context.Context, path string, args ...any) error {
	queryString, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Unable to find sql for creating posts table: %v\n", err)
		os.Exit(1)
	}
	_, err = app.Db.Exec(ctx, string(queryString), args...)
	return err
}

func FetchRows[T any](app *App, ctx context.Context, path string, args ...any) ([]T, error) {
	queryString, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Unable to find sql for creating posts table: %v\n", err)
		os.Exit(1)
	}
	rows, err := app.Db.Query(ctx, string(queryString), args...)
	return pgx.CollectRows[T](rows, pgx.RowToStructByName[T])
}

func (app *App) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "app", app)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FetchRow[T any](app *App, ctx context.Context, path string, args ...any) (*T, error) {
	queryString, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to find sql for creating posts table: %v\n", err)
	}
	rows, err := app.Db.Query(ctx, string(queryString), args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[T])
}
