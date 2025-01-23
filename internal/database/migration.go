package database

import (
	"log"

	"github.com/munene-m/pepa/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    // Create tables in the correct order
    err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.AutoMigrate(&models.Invoice{}); err != nil {
            return err
        }
        
        if err := tx.AutoMigrate(&models.Item{}); err != nil {
            return err
        }
        if err := tx.AutoMigrate(&models.User{}); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        log.Printf("Migration failed: %v", err)
        return err
    }
    
    return nil
}