package middlewares 

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"

)

func FirebaseAuthMiddleware(authClient *auth.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Look for the "session" cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			// If no cookie or error, send 401 and redirect via HTMX header
			// Log the failed authentication, but don't expose error to client
			w.Header().Set("HX-Redirect", "/login") // HTMX-specific redirect
			http.Error(w, "Session cookie required", http.StatusUnauthorized)
			return
		}

		// 2. Verify the session cookie using the Admin SDK
		// This checks validity, expiration, and ensures it hasn't been revoked.
		token, err := authClient.VerifySessionCookie(r.Context(), cookie.Value)
		if err != nil {
			// If invalid, clear the cookie and redirect
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})
			w.Header().Set("HX-Redirect", "/login")
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// 3. Cookie is valid: set the user's UID in the request context
		ctx := context.WithValue(r.Context(), "uid", token.UID)
		user, err := authClient.GetUser(ctx, token.UID)
		ctx = context.WithValue(ctx, "user", user)

		// 4. Proceed to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
