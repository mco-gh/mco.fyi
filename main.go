package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	fmt.Println(path)
	if path == "" {
		http.ServeFile(w, r, "home.html")
	} else if path == "meiko.jpg" {
		http.ServeFile(w, r, "meiko.jpg")
	} else if url, ok := links[path]; ok {
		//do something here
		http.Redirect(w, r, url.(string), 301)
	} else {
		http.ServeFile(w, r, "404.html")
	}
}

var links map[string]interface{}

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
	links = docsnap.Data()

	fmt.Println(links)

	http.HandleFunc("/", redirect)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
