package database

import (
	"database/sql"
	"log"

	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeConnection(dbChannel chan *gorm.DB, sqldbChannel chan *sql.DB) {
	dsn := "root:root@tcp(0.0.0.0:3306)/invoices?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	log.Println("Connected to Database!")

	sdb, errDBGet := db.DB()
	if errDBGet != nil {
		log.Println(errDBGet)
	}
	errPing := sdb.Ping()
	if errPing != nil {
		panic(errPing)
	}
	setupDatabase(db)
	dbChannel <- db
	sqldbChannel <- sdb
}

func setupDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&dataStructures.Invoice{})
	if err != nil {
		log.Println(err)
	}
}
