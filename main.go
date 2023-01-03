package main

import (
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsawlen/matchingAppInvoiceService/common/database"
	"github.com/tsawlen/matchingAppInvoiceService/common/initializer"
	"github.com/tsawlen/matchingAppInvoiceService/controller"

	"github.com/google/uuid"
	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
	"gorm.io/gorm"
)

var mockCity = dataStructures.City{
	PLZ:   51147,
	Place: "Köln",
}

var mockUser = dataStructures.User{
	First_name:  "Erika",
	Name:        "Mustermann",
	Street:      "Heidestrasse",
	HouseNumber: "17",
	City:        &mockCity,
}

var mockInvoice = dataStructures.Invoice{
	CreatedAt: time.Now(),
	Id:        uuid.MustParse("c289c950-76fc-11ed-a1eb-0242ac120002"),
	Amount:    107.10,
	Service:   "Dämonische Beschwärungskunst Advanced Level",
	Hours:     2,
}

func main() {
	readyChannel := make(chan bool)
	go initializer.LoadEnvVariables(readyChannel)
	<-readyChannel
	dbChannel := make(chan *gorm.DB)
	sqlDBChannel := make(chan *sql.DB)
	go database.InitializeConnection(dbChannel, sqlDBChannel)

	db := <-dbChannel
	sqlDB := <-sqlDBChannel
	defer sqlDB.Close()

	router := gin.Default()

	router.GET("/invoice", controller.GetAllInvoices(db))
	router.GET("/invoice/user/:id", controller.GetAllInvoicesByUser(db))
	router.GET("/create-checkout-session", controller.CreateCheckoutSession())
	router.GET("/create-payment-intent", controller.CreatePaymentIntent())

	router.PUT("/invoice", controller.CreateInvoice(db))
	router.PUT("/invoice/pay/:id", controller.SetInvoiceToPayed(db))

	router.Run("0.0.0.0:8085")

	// STRIPE SECTION

}
