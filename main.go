package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var links map[string]interface{}

func redirect(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	fmt.Println(path)
	if path == "" || path == "/" {
		t, err := template.ParseFiles("home.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		err = t.Execute(w, links)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else if strings.HasPrefix(path, "css/") ||
		strings.HasPrefix(path, "img/") {
		http.ServeFile(w, r, path)
	} else if url, ok := links[path]; ok {
		http.Redirect(w, r, url.(string), 301)
	} else {
		http.ServeFile(w, r, "404.html")
	}
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
	links = docsnap.Data()
	http.HandleFunc("/", redirect)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
