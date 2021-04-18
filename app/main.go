package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	HelloWorldResponse struct {
		Message string `json:"message"`
	}
	HealthResponse struct {
		Status string `json:"status"`
	}
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/hello-world", helloWorldHandler)
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Listening on port 10000....")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request){
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

func healthHandler(w http.ResponseWriter, r *http.Request){
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