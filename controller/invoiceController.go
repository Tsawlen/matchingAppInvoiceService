package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
	"github.com/tsawlen/matchingAppInvoiceService/common/dbInterface"
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
		// Create Invoice
		invoice, errCreate := dbInterface.CreateInvoice(db, &newInvoice)
		if errCreate != nil {
			// To-Do: Publish to logger
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errCreate,
			})
			return
		}
		// To-Do: Generate Invoice PDF
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
