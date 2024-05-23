package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"RestaurantMenuService/prisma/db"
	
)

type AddItem struct {
	Name string   `json:"name"`
	Price float   `json:"price"`
	Quantity int   `json:"quantity"`
}


func (additem *AddItem) FromJSON(r io.Reader) error {
	d:= json.NewDecoder(r)
	return d.Decode(additem)
}


func POST_AddItem_Handler (w http.ResponseWriter, r *http.Request) {
	

	client := db.NewClient() 
	ctx := context.Background()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Printf("Error connecting database: %s",err.Error())
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			fmt.Printf("Error Disconnecting database: %s",err.Error())
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	
	var requestData AddItem
	if err := requestData.FromJSON(r.Body); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	
	_, err := client.Items.CreateOne(
		db.Items.Name.Set(requestData.Name),
		db.Items.Price.Set(requestData.Price),
		db.Items.Quantity.Set(requestData.Quantity),
	).Exec(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}