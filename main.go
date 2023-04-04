package main

import (
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func main() {
	// This is your test secret API key.
	stripe.Key = "PAYMENT_KEY"

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/create-checkout-session", createCheckoutSession)
	addr := "localhost:4242"
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	domain := "http://localhost:4242"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Fare"),
					},
					// $1.5 => 150
					UnitAmount: stripe.Int64(100),
				},
				Quantity: stripe.Int64(1),
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/success.html"),
		CancelURL:  stripe.String(domain + "/cancel.html"),
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}
