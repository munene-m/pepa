package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

	mux := http.NewServeMux()
	// Initialize handlers
	invoiceHandler := &handlers.InvoiceHandler{
		DB: db,
	}
	// Set up routes
	mux.HandleFunc("/api/invoices/create", invoiceHandler.CreateInvoice)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Welcome to the Pepa API!"))
	})

	// Use the mux as the main handler
	http.Handle("/", mux)

	log.Println("Server starting on port 9000...")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
