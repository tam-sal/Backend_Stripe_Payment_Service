package main

import (
	"backend/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type application struct {
	handlers Handlers
	config   Config
}

type Handlers struct {
}
type Config struct {
	stripeKey string
}

func (h Handlers) handleHealthCheck(w http.ResponseWriter, req *http.Request) {
	var response []byte = []byte("Server is up and running")
	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}
func calculateOrderAmount(productId string) int64 {
	switch productId {
	case "Forever Pants":
		return 26000
	case "Forever Shirt":
		return 15500
	case "Forever Shorts":
		return 30000
	}
	return 0
}
func (h Handlers) handleCreatePaymentIntent(w http.ResponseWriter, req *http.Request) {

	var request struct {
		ProductId string `json:"product_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address1  string `json:"address_1"`
		Address2  string `json:"address_2"`
		City      string `json:"city"`
		State     string `json:"state"`
		Zip       string `json:"zip"`
		Country   string `json:"country"`
	}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(request.ProductId)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var response struct {
		ClientSecret string `json:"clientSecret"`
	}
	response.ClientSecret = paymentIntent.ClientSecret
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, &buf)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	stripeK := env.GetString("STRIPE_KEY", "888")
	fmt.Println(stripeK)
	app := &application{
		handlers: Handlers{},
		config:   Config{stripeKey: stripeK},
	}
	stripe.Key = app.config.stripeKey
	http.HandleFunc("GET /health-check", app.handlers.handleHealthCheck)
	http.HandleFunc("POST /create-payment-intent", app.handlers.handleCreatePaymentIntent)

	fmt.Println("Server is listening on port 4242")
	err := http.ListenAndServe("0.0.0.0:4242", nil)
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}
