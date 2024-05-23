package main

import (
	"fmt"
	"net/http"
	"RestaurantMenuService/handlers"
	"RestaurantMenuService/kafka"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/additem", handlers.POST_AddItem_Handler).Methods("POST")
	r.HandleFunc("/items", handlers.GET_GetItems_Handler).Methods("GET")
	
	
	
	fmt.Println("Server is running...")
	err := http.ListenAndServe(":9000", r)
	if err!=nil{
		fmt.Printf("Cannot start server: %s",err.Error())
	}
}