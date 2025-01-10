package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/munene-m/pepa/internal/database"
	"github.com/munene-m/pepa/internal/handlers"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construct database connection string
	dsn := "host=localhost user=" + os.Getenv("DB_USER") + 
		" password=" + os.Getenv("DB_PASSWORD") + 
		" dbname=" + os.Getenv("DB_NAME") + 
		" port=" + os.Getenv("DB_PORT") + 
		" sslmode=disable"

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")
	handlers.InitializeGoogleAuth()

	if err := database.Migrate(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

	mux := http.NewServeMux()
	// Initialize handlers
	invoiceHandler := &handlers.InvoiceHandler{
		DB: db,
	}
	authHandler := handlers.NewAuthHandler(db)
	// Set up routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/auth/google/signin", authHandler.SignInWithProvider)
	mux.HandleFunc("/auth/google/callback", authHandler.CallbackHandler)
	mux.HandleFunc("/auth/logout", authHandler.LogoutHandler)
	mux.HandleFunc("/success", authHandler.RequireAuth(authHandler.Success))

	// Protected routes
	mux.HandleFunc("/api/invoices/create", authHandler.RequireAuth(invoiceHandler.CreateInvoice))

	// Use the mux as the main handler
	http.Handle("/", mux)

	log.Println("Server starting on port 9000...")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
