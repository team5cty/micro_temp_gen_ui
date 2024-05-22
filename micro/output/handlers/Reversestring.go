package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"ReverseService/prisma/db"
	
)

type Reversestring struct {
	OriginalString string   `json:"originalstring"`
}


func (reversestring *Reversestring) FromJSON(r io.Reader) error {
	d:= json.NewDecoder(r)
	return d.Decode(reversestring)
}


func POST_Reversestring_Handler (w http.ResponseWriter, r *http.Request) {
	

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
	
	var requestData Reversestring
	if err := requestData.FromJSON(r.Body); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	
	_, err := client.ReversedStrings.CreateOne(
		db.ReversedStrings.OriginalString.Set(requestData.OriginalString),
	).Exec(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}