package dbInterface

import (
	"time"

	"github.com/google/uuid"
	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
	"gorm.io/gorm"
)

func GetAllInvoices(db *gorm.DB) (*[]dataStructures.Invoice, error) {
	var invoices *[]dataStructures.Invoice
	result := db.Find(&invoices)
	if result.Error != nil {
		return nil, result.Error
	}
	return invoices, nil
}

func GetAllInvoicesByUser(db *gorm.DB, id string) (*[]dataStructures.Invoice, error) {
	var invoices *[]dataStructures.Invoice

	err := db.Model(&dataStructures.Invoice{}).Where("payer=?", id).Find(&invoices).Error

	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func CreateInvoice(db *gorm.DB, invoice *dataStructures.Invoice) (*dataStructures.Invoice, error) {
	result := db.Create(&invoice)
	if result.Error != nil {
		return &dataStructures.Invoice{}, result.Error
	}
	return invoice, nil
}

func SetInvoiceToPayed(db *gorm.DB, invoiceID *uuid.UUID) (bool, error) {
	now := time.Now().Local()
	result := db.Model(&dataStructures.Invoice{}).Where("id = ?", invoiceID).Update("payed_at", now)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
