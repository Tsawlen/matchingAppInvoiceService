package publish

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
)

type invoiceMessage struct {
	ToUser  uint
	Type    string
	Message []byte
}

func PublishInvoice(invoice *dataStructures.Invoice) {
	client := &http.Client{}
	jsonInv, errInv := json.Marshal(invoice)
	if errInv != nil {
		log.Println(errInv)
	}
	messageObject := invoiceMessage{
		ToUser:  uint(invoice.Payer),
		Type:    "invoice",
		Message: jsonInv,
	}
	json, errMes := json.Marshal(messageObject)
	if errMes != nil {
		log.Println(errMes)
	}
	req, errReq := http.NewRequest(http.MethodPut, os.Getenv("CLOUD_RELAY_PUB"), bytes.NewBuffer(json))
	if errReq != nil {
		log.Println(errReq)
	}
	req.Header.Set("topic", "email")
	req.Header.Set("service", "Invoice Service")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, errSend := client.Do(req)
	if errSend != nil {
		log.Println(errSend)
	}
}
