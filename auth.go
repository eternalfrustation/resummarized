package main

import (
	"log"
	"net/http"
	"time"
)

// Define the request body structure
type TokenRequest struct {
	IDToken string `json:"idToken"`
}

// SessionLoginHandler is where the Firebase ID Token is exchanged for a session cookie.
func SessionLoginHandler(w http.ResponseWriter, r *http.Request, tokenReq TokenRequest) {
	// 2. Set the desired cookie duration (e.g., 5 days)
	// Firebase recommends a maximum of 2 weeks (1209600 seconds)
	expiresIn := time.Hour * 24 * 5

	// 3. CRUCIAL: Create the secure session cookie using the Admin SDK
	sessionCookie, err := authClient.SessionCookie(r.Context(), tokenReq.IDToken, expiresIn)
	if err != nil {
		log.Printf("Failed to create session cookie: %v", err)
		http.Error(w, "Failed to create session", http.StatusUnauthorized)
		return
	}

	// 4. Set the HTTP-Only cookie in the browser
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		MaxAge:   int(expiresIn.Seconds()),
		HttpOnly: true, // Prevents client-side JS access (HIGHLY recommended)
		Secure:   true, // Requires HTTPS (HIGHLY recommended for production)
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Or Strict
	})

	w.Header().Set("HX-Redirect", "/dashboard")

	w.WriteHeader(200)
	w.Write([]byte("{'success': true}"))
}
