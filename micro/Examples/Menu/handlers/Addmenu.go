package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"Menu/prisma/db"
	
)

type Addmenu struct {
	Availqty int   `json:"availqty"`
	Desc string   `json:"desc"`
	Menuid int   `json:"menuid"`
	Name string   `json:"name"`
}


func (addmenu *Addmenu) FromJSON(r io.Reader) error {
	d:= json.NewDecoder(r)
	return d.Decode(addmenu)
}


func POST_Addmenu_Handler (w http.ResponseWriter, r *http.Request) {
	

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
	
	var requestData Addmenu
	if err := requestData.FromJSON(r.Body); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	
	_, err := client.Menu.CreateOne(
		db.Menu.Availqty.Set(requestData.Availqty),
		db.Menu.Desc.Set(requestData.Desc),
		db.Menu.Menuid.Set(requestData.Menuid),
		db.Menu.Name.Set(requestData.Name),
	).Exec(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}