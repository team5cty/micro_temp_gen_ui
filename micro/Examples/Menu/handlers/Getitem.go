package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"strconv"

	"Menu/prisma/db"

	"github.com/gorilla/mux"
)

type Getitem struct {
	Availqty int    `json:"availqty"`
	Desc     string `json:"desc"`
	Name     string `json:"name"`
}

func (getitem *Getitem) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(getitem)
}

func GET_Getitem_Handler(w http.ResponseWriter, r *http.Request) {

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
	m := mux.Vars(r)
	var val string
	for _, v := range m {
		val = v
	}
	value, _ := strconv.Atoi(val)

	res, err := client.Menu.FindUnique(db.Menu.Menuid.Equals(value)).Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ele := &Getitem{
		Availqty: res.Availqty,
		Desc:     res.Desc,
		Name:     res.Name,
	}
	ele.ToJSON(w)
}
