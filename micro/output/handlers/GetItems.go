package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"RestaurantMenuService/prisma/db"
	
)

type GetItems struct {
	Itemid Int   `json:"itemid"`
	Name string   `json:"name"`
	Price float   `json:"price"`
	Quantity int   `json:"quantity"`
}
type GetItems_list []*GetItems


func (getitems *GetItems_list) ToJSON(w io.Writer) error {
	e:= json.NewEncoder(w)
	return e.Encode(getitems)
}


func GET_GetItems_Handler (w http.ResponseWriter, r *http.Request) {
	

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
	var getitems GetItems_list
	res, err := client.Items.FindMany().Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, object := range res {
		ele := &GetItems{
			Itemid: object.Itemid,
			Name: object.Name,
			Price: object.Price,
			Quantity: object.Quantity,
		}
		getitems = append(getitems, ele)
	}
	getitems.ToJSON(w)
}