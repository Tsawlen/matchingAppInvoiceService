package controller

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func CreatePaymentIntent() gin.HandlerFunc {
	handler := func(context *gin.Context) {
		stripe.Key = os.Getenv("STRIPE_KEY")
		price := context.Request.Header.Get("Price")

		priceInt, errConv := strconv.ParseFloat(price, 64)
		if errConv != nil {
			log.Println(errConv)
			return
		}
		finishedPrice := (int64(priceInt * 100))

		params := &stripe.PaymentIntentParams{
			Amount:   &finishedPrice,
			Currency: stripe.String(string(stripe.CurrencyEUR)),
		}
		pi, _ := paymentintent.New(params)
		context.IndentedJSON(http.StatusOK, pi)
	}
	return gin.HandlerFunc(handler)
}

func CreateCheckoutSession() gin.HandlerFunc {
	handler := func(context *gin.Context) {
		stripe.Key = os.Getenv("STRIPE_KEY")
		price := context.Request.Header.Get("Price")
		service := context.Request.Header.Get("Service")

		priceInt, errConv := strconv.ParseFloat(price, 64)
		description := "Payment for Finder"
		if errConv != nil {
			log.Println(errConv)
			return
		}
		finishedPrice := (int64(priceInt * 100))
		domain := "http://localhost:8090"
		params := &stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				&stripe.CheckoutSessionLineItemParams{
					PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
						UnitAmount: &finishedPrice,
						Currency:   stripe.String(string(stripe.CurrencyEUR)),
						ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
							Name:        &service,
							Description: &description,
						},
					},
					Quantity: stripe.Int64(1),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(domain + "?success=true"),
			CancelURL:  stripe.String(domain + "?canceled=true"),
		}

		s, err := session.New(params)
		if err != nil {
			log.Printf("session.new: %v", err)
		}
		context.Redirect(http.StatusSeeOther, s.URL)
	}
	return gin.HandlerFunc(handler)
}

func CreateCheckOutSession(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	price := r.Header.Get("Price")
	service := r.Header.Get("Service")
	priceInt, errConv := strconv.ParseFloat(price, 64)
	description := "Payment for Finder"
	if errConv != nil {
		log.Println(errConv)
		return
	}
	finishedPrice := (int64(priceInt * 100))
	domain := "http://localhost:8090"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: &finishedPrice,
					Currency:   stripe.String(string(stripe.CurrencyEUR)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        &service,
						Description: &description,
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "?success=true"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("session.new: %v", err)
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}
