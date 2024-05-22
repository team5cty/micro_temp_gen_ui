package main

import (
	"Menu/handlers"
	"Menu/kafka"
	"Menu/prisma/db"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/addmenu", handlers.POST_Addmenu_Handler).Methods("POST")
	r.HandleFunc("/menu", handlers.GET_Getmenu_Handler).Methods("GET")
	r.HandleFunc("/menu/{id}", handlers.GET_Getitem_Handler).Methods("GET")

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

	go kafka.Consume("quantity", 0, func(s string) {
		split := strings.Split(s, ",")
		pids := split[1]
		qtys := split[0]
		pid, _ := strconv.Atoi(pids)
		qty, _ := strconv.Atoi(qtys)

		menu, _ := client.Menu.FindUnique(
			db.Menu.Menuid.Equals(pid),
		).Exec(ctx)
		newqty := menu.Availqty - qty
		client.Menu.FindUnique(
			db.Menu.Menuid.Equals(pid),
		).Update(
			db.Menu.Availqty.Set(newqty),
		).Exec(ctx)
	})

	fmt.Println("Server is running...")
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		fmt.Printf("Cannot start server: %s", err.Error())
	}
}
