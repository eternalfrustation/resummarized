package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"
	"time"

	"crypto/rand"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/UniquityVentures/resummarized/layouts"
	"github.com/UniquityVentures/resummarized/pages"
	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rickb777/servefiles/v3"
	"google.golang.org/api/option"
)

var decoder = schema.NewDecoder()
var authClient *auth.Client // Global client for re-use

func initFirebase() {
	// 1. Point to your Service Account JSON file path
	// IMPORTANT: Keep this file secure and do not commit it to source control.
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

type App struct {
	Db *pgx.Conn
}

func DbConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

type ArticleSummary struct {
	ArticleID int `json:"article_id" comment:"Database ID for internal tracking. 0 when c"`

	HeadlineTitle string `json:"headline_title" comment:"The catchy, non-technical title/headline."`
	LeadParagraph string `json:"lead_paragraph" comment:"The opening paragraph that states the main finding and its real-world impact (the 'Why You Care')."`

	BackgroundContext string `json:"background_context" comment:"Simplified background on the large problem or field of study."`
	ResearchQuestion  string `json:"research_question" comment:"The focused question the specific research paper set out to answer, framed as a mystery or challenge."`

	SimplifiedMethods string `json:"simplified_methods" comment:"A non-expert friendly description of the methodology (what they did, not how the instruments work)."`

	CoreFindings    string `json:"core_findings" comment:"The primary result/conclusion, stated simply and focused on magnitude or effect."`
	SurpriseFinding string `json:"surprise_finding,omitempty" comment:"Any unexpected or counter-intuitive results found (optional)."`

	FutureImplications string `json:"future_implications" comment:"What the discovery means for humanity, technology, or the field (the future potential)."`
	StudyLimitations   string `json:"study_limitations" comment:"The most critical limitation or caveat of the study (e.g., 'only conducted in mice')."`
	NextSteps          string `json:"next_steps" comment:"The immediate next phase of research or development."`

	Author        string    `json:"author" comment:"The name of the article writer (not necessarily the scientist)."`
	DatePublished time.Time `json:"date_published" comment:"Date the summary article was published."`
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

func main() {
	initFirebase()
	app := initApp()
	csrfKey := make([]byte, 128)
	rand.Read(csrfKey)
	csrfMiddleware := csrf.Protect(csrfKey, csrf.TrustedOrigins([]string{"localhost:7331", "localhost:4269"}), csrf.FieldName("_csrf"))
	_ = csrfMiddleware

	router := http.NewServeMux()
	router.Handle("/layout_test", templ.Handler(layouts.Navbar()))
	router.Handle("/login", (http.HandlerFunc(LoginHandler)))
	router.Handle("/dashboard", FirebaseAuthMiddleware(authClient, http.HandlerFunc(DashboardHandler)))
	router.Handle("/create", FirebaseAuthMiddleware(authClient, http.HandlerFunc(CreateHandler)))
	router.Handle("/assets/", servefiles.NewAssetHandler("./assets/").StripOff(1) /*.WithMaxAge(24 * time.Hour) */)
	router.HandleFunc("/auth/sessionLogin/", SessionLoginHandler)
	decoder.IgnoreUnknownKeys(true)

	log.Println("Listening on :4269")
	http.ListenAndServe(":4269", app.Handle(router))
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	templ.Handler(pages.Dashboard()).ServeHTTP(w, r)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Invalid Submission"))
		}
		var postCreateForm pages.PostCreateForm

		err = decoder.Decode(&postCreateForm, r.PostForm)
		fmt.Println(postCreateForm)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		app := r.Context().Value("app").(*App)

		create_post, err := os.ReadFile("sql/create_post.sql")
		if err != nil {
			log.Printf("Unable to find sql for creating posts table: %v\n", err)
			os.Exit(1)
		}
		app.Db.QueryRow(r.Context(), string(create_post),

			postCreateForm.HeadlineTitle,      // $1
			postCreateForm.LeadParagraph,      // $2
			postCreateForm.BackgroundContext,  // $3
			postCreateForm.ResearchQuestion,   // $4
			postCreateForm.SimplifiedMethods,  // $5
			postCreateForm.CoreFindings,       // $6
			postCreateForm.SurpriseFindings,   // $7 (If empty string is passed, it stores "" not NULL)
			postCreateForm.FutureImplications, // $8
			postCreateForm.StudyLimitations,   // $9
			postCreateForm.NextSteps,          // $10
		)
		w.Header().Set("HX-Redirect", "/dashboard")
	}

	templ.Handler(pages.PostCreatePage()).ServeHTTP(w, r)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Invalid Submission"))
	}
	var loginForm pages.LoginForm

	err = decoder.Decode(&loginForm, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templ.Handler(pages.Login()).ServeHTTP(w, r)
}

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

// Define the request body structure
type tokenRequest struct {
	IDToken string `json:"idToken"`
}

// SessionLoginHandler is where the Firebase ID Token is exchanged for a session cookie.
func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decode the request body to get the ID Token
	var req tokenRequest

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.IDToken = r.FormValue("idToken")

	// 2. Set the desired cookie duration (e.g., 5 days)
	// Firebase recommends a maximum of 2 weeks (1209600 seconds)
	expiresIn := time.Hour * 24 * 5

	// 3. CRUCIAL: Create the secure session cookie using the Admin SDK
	sessionCookie, err := authClient.SessionCookie(r.Context(), req.IDToken, expiresIn)
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
