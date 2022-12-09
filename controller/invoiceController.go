package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	publish "github.com/tsawlen/matchingAppInvoiceService/Publish"
	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
	"github.com/tsawlen/matchingAppInvoiceService/common/dbInterface"
	"github.com/tsawlen/matchingAppInvoiceService/connector"
	"github.com/tsawlen/matchingAppInvoiceService/documentGenerator"
	"gorm.io/gorm"
)

func GetAllInvoices(db *gorm.DB) gin.HandlerFunc {
	handler := func(context *gin.Context) {
		invoices, err := dbInterface.GetAllInvoices(db)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}
		context.JSON(http.StatusOK, invoices)
	}
	return gin.HandlerFunc(handler)
}

func GetAllInvoicesByUser(db *gorm.DB) gin.HandlerFunc {
	handler := func(context *gin.Context) {
		id := context.Param("id")
		invoices, err := dbInterface.GetAllInvoicesByUser(db, id)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}
		context.JSON(http.StatusOK, invoices)
	}
	return gin.HandlerFunc(handler)
}

func CreateInvoice(db *gorm.DB) gin.HandlerFunc {
	handler := func(context *gin.Context) {
		var newInvoice dataStructures.Invoice
		if err := context.BindJSON(&newInvoice); err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		// To-Do: Check for booth IDs to be valid
		billed, errBill := connector.GetProfileById(newInvoice.Payer)
		if errBill != nil || billed.ID == 0 {
			log.Println(errBill)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errBill,
			})
			return
		}
		biller, errBiller := connector.GetProfileById(newInvoice.Biller)
		if errBiller != nil || biller.ID == 0 {
			log.Println(errBill)
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errBill,
			})
			return
		}
		// Attributes
		uuid, errUUID := uuid.NewRandom()
		newInvoice.Id = uuid
		if errUUID != nil {
			log.Println(errUUID)
		}
		newInvoice.Amount = float64(newInvoice.Hours) * biller.Price
		// Create InvoicePDF
		pdfInvoice := documentGenerator.GenerateInvoicePDF(billed, &newInvoice)
		newInvoice.InvoicePDF = pdfInvoice
		// Create Invoice
		invoice, errCreate := dbInterface.CreateInvoice(db, &newInvoice)
		if errCreate != nil {
			// To-Do: Publish to logger
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errCreate,
			})
			return
		}
		publish.PublishInvoice(invoice)
		context.JSON(http.StatusCreated, invoice)
	}
	return gin.HandlerFunc(handler)
}

func SetInvoiceToPayed(db *gorm.DB) gin.HandlerFunc {
	handler := func(context *gin.Context) {
		id := context.Param("id")
		if id == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "No UUID given to pay",
			})
			return
		}
		// Parse UUID
		invoiceUuid, parseErr := uuid.Parse(id)
		if parseErr != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": parseErr,
			})
			return
		}
		ok, err := dbInterface.SetInvoiceToPayed(db, &invoiceUuid)
		if !ok {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"message": "Invoice set to payed!",
		})
		return
	}
	return gin.HandlerFunc(handler)
}
