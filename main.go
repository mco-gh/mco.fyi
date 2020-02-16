package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"net/http"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.google.com", 301)
}

func main() {
	proj := "mco-fyi"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, proj)
	if err != nil {
		log.Fatalln(err)
	}
	shortlinks := client.Doc("Redirects/Shortlinks")

	docsnap, err := shortlinks.Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	dataMap := docsnap.Data()
	fmt.Println(dataMap)

	http.HandleFunc("/", redirect)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
