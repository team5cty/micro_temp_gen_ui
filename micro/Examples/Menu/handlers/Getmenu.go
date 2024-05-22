package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"Menu/prisma/db"
)

type Getmenu struct {
	Desc string `json:"desc"`
	Name string `json:"name"`
}
type Getmenu_list []*Getmenu

func (getmenu *Getmenu_list) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(getmenu)
}

func GET_Getmenu_Handler(w http.ResponseWriter, r *http.Request) {

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
	var getmenu Getmenu_list
	res, err := client.Menu.FindMany().Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, object := range res {
		ele := &Getmenu{
			Desc: object.Desc,
			Name: object.Name,
		}
		getmenu = append(getmenu, ele)
	}
	getmenu.ToJSON(w)
}
