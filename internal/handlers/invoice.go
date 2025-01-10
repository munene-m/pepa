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
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
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

	var items []models.Item
    itemDescriptions := r.Form["items.item_description"]
    itemQuantities := r.Form["items.item_quantity"]
    itemRates := r.Form["items.item_rate"]
    
    for i := 0; i < len(itemDescriptions); i++ {
        items = append(items, models.Item{
            ItemDescription: itemDescriptions[i],
            ItemQuantity:    parseInt(itemQuantities[i]),
            ItemRate:        parseFloat(itemRates[i]),
            Amount:          parseFloat(itemRates[i]) * float64(parseInt(itemQuantities[i])),
        })
    }

	// Parse invoice number to int64
	invoiceNumber, err := strconv.ParseInt(r.FormValue("invoice_number"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid invoice number: "+err.Error(), http.StatusBadRequest)
		return
	}

	invoice := models.Invoice{
		BillTo:         r.FormValue("bill_to"),
		BillFrom:       r.FormValue("bill_from"), 
		InvoiceNumber: invoiceNumber,
		InvoiceDate:    r.FormValue("invoice_date"),
		DueDate:        r.FormValue("due_date"),
		Items:          items,
		InvoiceSubTotal: parseFloat(r.FormValue("invoice_sub_total")),
		Tax:            parseFloat(r.FormValue("tax")),
		InvoiceTotal:   parseFloat(r.FormValue("invoice_total")),
		InvoiceStatus:  r.FormValue("invoice_status"),
		Notes:          r.FormValue("notes"),
		Currency:       r.FormValue("currency"),
		BillToLogo:     logoPath,
	}
	transactionErr := h.DB.Transaction(func(tx *gorm.DB) error {
        // Create the invoice first
        if err := tx.Create(&invoice).Error; err != nil {
            return err
        }

        // Create associated items
        for i := range items {
            items[i].InvoiceID = invoice.ID
            if err := tx.Create(&items[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })

	if transactionErr != nil {
		http.Error(w, "Failed to create invoice: "+transactionErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}