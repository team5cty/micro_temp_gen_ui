package main

import (
	"Order/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/order", handlers.POST_Order_Handler).Methods("POST")
	r.HandleFunc("/details/{id}", handlers.GET_Order_details_Handler).Methods("GET")

	fmt.Println("Server is running...")
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		fmt.Printf("Cannot start server: %s", err.Error())
	}
}
