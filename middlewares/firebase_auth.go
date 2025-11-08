package middlewares

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/UniquityVentures/resummarized/core"
)

func FirebaseAuthMiddleware(authClient *auth.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := authClient.VerifySessionCookie(r.Context(), cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "uid", token.UID)
		user, err := authClient.GetUser(ctx, token.UID)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("uid").(string)
		if !ok {
			w.Header().Set("HX-Redirect", "/login") // HTMX-specific redirect
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}

		_, ok = r.Context().Value("app").(*core.App)
		if !ok {
			w.Header().Set("HX-Redirect", "/login") // HTMX-specific redirect
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("uid").(string)
		if !ok {
			w.Header().Set("HX-Redirect", "/login") // HTMX-specific redirect
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}

		app, ok := r.Context().Value("app").(*core.App)
		if !ok {
			w.Header().Set("HX-Redirect", "/login") // HTMX-specific redirect
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !ok {
			log.Fatalln("App not inject into requests")
		}
		if isAdmin, err := app.IsAdmin(r.Context(), uid); !isAdmin {
			w.Header().Set("HX-Redirect", "/login")
			http.Error(w, "Not admin", http.StatusUnauthorized)
			if err != nil {
				log.Println(err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
