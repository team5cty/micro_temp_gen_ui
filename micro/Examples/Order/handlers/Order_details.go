package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"strconv"
	
	"github.com/gorilla/mux"
	"Order/prisma/db"
	
)

type Order_details struct {
	Productid int   `json:"productid"`
	Quantity int   `json:"quantity"`
	Status string   `json:"status"`
}


func (order_details *Order_details) ToJSON(w io.Writer) error {
	e:= json.NewEncoder(w)
	return e.Encode(order_details)
}


func GET_Order_details_Handler (w http.ResponseWriter, r *http.Request) {
	

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
	m:=mux.Vars(r)
	var val string
	for _,v := range m{
		val=v
	}
	value, _ := strconv.Atoi(val)

	res, err := client.Order.FindUnique(db.Order.Orderid.Equals(value)).Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ele := &Order_details{
		Productid:res.Productid,
		Quantity:res.Quantity,
		Status:res.Status,
	}
	ele.ToJSON(w)
}