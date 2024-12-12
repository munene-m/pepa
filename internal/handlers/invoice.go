package handlers

import (
	"context"
	"encoding/json"

	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/munene-m/pepa/internal/models"
	"github.com/munene-m/pepa/internal/services"
)

type InvoiceHandler struct {
    DB *gorm.DB
}

// Helper function for parsing float values
func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Set max file size to 10MB
	r.ParseMultipartForm(10 << 20)

	 // Get the file from form data
	 file, _, err := r.FormFile("bill_to_logo")
	 if err != nil && err != http.ErrMissingFile {
		 http.Error(w, "Error retrieving logo file: "+err.Error(), http.StatusBadRequest)
		 return
	 }
 
	 var logoPath string
	 if file != nil {
        defer file.Close()
        uploader, err := services.NewFileUploader()
        if err != nil {
            http.Error(w, "Failed to initialize file uploader: "+err.Error(), http.StatusInternalServerError)
            return
        }

        ctx := context.Background()
        err = uploader.CreateBucketIfNotExists(ctx, "us-east-1")
        if err != nil {
            http.Error(w, "Failed to create bucket: "+err.Error(), http.StatusInternalServerError)
            return
        }

        logoPath, err = uploader.UploadFile(ctx, file, "logos")
        if err != nil {
            http.Error(w, "Failed to upload logo: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

	invoice := models.Invoice{
		BillTo:         r.FormValue("bill_to"),
		BillFrom:       r.FormValue("bill_from"), 
		InvoiceNumber:  r.FormValue("invoice_number"),
		InvoiceDate:    r.FormValue("invoice_date"),
		DueDate:        r.FormValue("due_date"),
		InvoiceSubTotal: parseFloat(r.FormValue("invoice_sub_total")),
		Tax:            parseFloat(r.FormValue("tax")),
		InvoiceTotal:   parseFloat(r.FormValue("invoice_total")),
		InvoiceStatus:  r.FormValue("invoice_status"),
		Notes:          r.FormValue("notes"),
		Currency:       r.FormValue("currency"),
		BillToLogo:     logoPath,
	}

	result := h.DB.Create(&invoice)
	if result.Error != nil {
		http.Error(w, "Failed to create invoice: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}