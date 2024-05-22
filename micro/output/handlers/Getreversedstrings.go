package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"context"
	"fmt"
	
	"ReverseService/prisma/db"
	
)

type Getreversedstrings struct {
	OriginalString string   `json:"originalstring"`
	ReversedString string   `json:"reversedstring"`
}
type Getreversedstrings_list []*Getreversedstrings


func (getreversedstrings *Getreversedstrings_list) ToJSON(w io.Writer) error {
	e:= json.NewEncoder(w)
	return e.Encode(getreversedstrings)
}


func GET_Getreversedstrings_Handler (w http.ResponseWriter, r *http.Request) {
	

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
	var getreversedstrings Getreversedstrings_list
	res, err := client.ReversedStrings.FindMany().Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, object := range res {
		ele := &Getreversedstrings{
			OriginalString: object.OriginalString,
			ReversedString: object.ReversedString,
		}
		getreversedstrings = append(getreversedstrings, ele)
	}
	getreversedstrings.ToJSON(w)
}