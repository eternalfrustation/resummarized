package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var authClient *auth.Client // Global client for re-use

func initFirebase() {
	serviceAccountFile := os.Getenv("FIREBASE_SA_PATH")
	if serviceAccountFile == "" {
		log.Fatal("FIREBASE_SA_PATH environment variable not set.")
	}

	// 2. Initialize the Firebase App
	opt := option.WithCredentialsFile(serviceAccountFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	// 3. Get the Auth Client (this is what verifies tokens)
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v\n", err)
	}
	log.Println("Firebase Admin SDK initialized successfully.")
}
