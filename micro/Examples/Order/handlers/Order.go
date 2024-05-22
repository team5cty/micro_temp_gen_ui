package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"Order/prisma/db"

	"Order/kafka"
)

type Order struct {
	Orderid   int    `json:"orderid"`
	Productid int    `json:"productid"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

func (order *Order) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(order)
}

func POST_Order_Handler(w http.ResponseWriter, r *http.Request) {

	produce := kafka.Producer("quantity", 0)

	client := db.NewClient()
	ctx := context.Background()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Printf("Error connecting database: %s", err.Error())
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			fmt.Printf("Error Disconnecting database: %s", err.Error())
		}
	}()

	w.Header().Set("Content-Type", "application/json")

	var requestData Order
	if err := requestData.FromJSON(r.Body); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	_, err := client.Order.CreateOne(
		db.Order.Orderid.Set(requestData.Orderid),
		db.Order.Productid.Set(requestData.Productid),
		db.Order.Quantity.Set(requestData.Quantity),
		db.Order.Status.Set(requestData.Status),
	).Exec(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	produce(strconv.Itoa(requestData.Quantity) + "," + strconv.Itoa(requestData.Productid))
}
