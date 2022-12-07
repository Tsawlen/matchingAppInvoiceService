package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/tsawlen/matchingAppInvoiceService/common/database"
	"github.com/tsawlen/matchingAppInvoiceService/controller"
	"gorm.io/gorm"
)

func main() {

	dbChannel := make(chan *gorm.DB)
	sqlDBChannel := make(chan *sql.DB)
	go database.InitializeConnection(dbChannel, sqlDBChannel)

	db := <-dbChannel
	sqlDB := <-sqlDBChannel
	defer sqlDB.Close()

	router := gin.Default()

	router.GET("/invoice", controller.GetAllInvoices(db))

	router.PUT("/invoice", controller.CreateInvoice(db))
	router.PUT("/invoice/pay/:id", controller.SetInvoiceToPayed(db))

	router.Run("0.0.0.0:8085")

}
