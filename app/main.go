package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/rizkiduwinanto/duwinanto-golang-training-beginner/payment"
	"gopkg.in/go-playground/validator.v9"
)

type (
	HelloWorldResponse struct {
		Message string `json:"message"`
	}
	HealthResponse struct {
		Status string `json:"status"`
	}
)

type PaymentCodeHandler struct {
	service payment.Service
}

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/hello-world", helloWorldHandler)
	http.HandleFunc("/health", healthHandler)

	http.HandleFunc("/payment-codes", paymentCodeHandler)
	http.HandleFunc("/payment-codes/", paymentCodeHandler)

	fmt.Println("Listening on port 10000....")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	helloWorld := HelloWorldResponse{"hello world"}

	res, err := json.Marshal(helloWorld)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	fmt.Println(r.Method, r.URL.Path, http.StatusOK)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	health := HealthResponse{"healthy"}

	res, err := json.Marshal(health)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	fmt.Println(r.Method, r.URL.Path, http.StatusOK)
}

func (p PaymentCodeHandler) readBody(r *http.Request) (payment.Payment, error) {
	decoder := json.NewDecoder(r.Body)
	var paymentCode payment.Payment
	err := decoder.Decode(&paymentCode)
	if err != nil {
		return paymentCode, err
	}
	return paymentCode, nil
}

func (p PaymentCodeHandler) createPaymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	paymentCode, err := p.readBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	err = validate.Struct(paymentCode)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				http.Error(w, fmt.Sprintf("field '%s' is required", err.Field()), http.StatusBadRequest)
				return
			}
		}
	}

	ctx := context.TODO()

	err = p.service.Create(ctx, &paymentCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(paymentCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (p PaymentCodeHandler) getPaymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payment-codes/")

	ctx := context.TODO()

	paymentCode, err := p.service.Get(ctx, id)
	if err != nil {
		if err == pg.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(paymentCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func NewPaymentCodeHandler(s payment.Service) PaymentCodeHandler {
	return PaymentCodeHandler{
		service: s,
	}
}

func paymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	db := connectDB()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	repository := payment.NewRepository(db)
	service := payment.NewService(repository)
	handler := NewPaymentCodeHandler(service)

	switch r.Method {
	case "POST":
		handler.createPaymentCodeHandler(w, r)
		return
	case "GET":
		handler.getPaymentCodeHandler(w, r)
		return
	default:
		http.NotFound(w, r)
	}
}

func connectDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "babel_group",
		Password: "babel_group",
		Database: "babel_database",
	})
	return db
}
