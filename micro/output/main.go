package main

import (
	"fmt"
	"net/http"
	"ReverseService/handlers"
	"ReverseService/kafka"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reversestring", handlers.POST_Reversestring_Handler).Methods("POST")
	r.HandleFunc("/reversedstrings", handlers.GET_Getreversedstrings_Handler).Methods("GET")
	
	
	
	fmt.Println("Server is running...")
	err := http.ListenAndServe(":9001", r)
	if err!=nil{
		fmt.Printf("Cannot start server: %s",err.Error())
	}
}