package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	uuid "github.com/gofrs/uuid"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var links map[string]interface{}
var gaPropertyID = "UA-158788691-1"

func redirect(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
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
		if err := trackEvent(r, "short-links", "redirect", path, nil); err != nil {
			fmt.Fprintf(w, "Event did not track: %v", err)
			return
		}
		http.Redirect(w, r, url.(string), 301)
	} else {
		http.ServeFile(w, r, "404.html")
	}
}

func trackEvent(r *http.Request, category, action, label string, value *uint) error {
	if category == "" || action == "" {
		return errors.New("analytics: category and action are required")
	}
	v := url.Values{
		"v":   {"1"},
		"tid": {gaPropertyID},
		// Anonymously identifies a particular user. See the parameter guide for
		// details:
		// https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid
		//
		// Depending on your application, this might want to be associated with the
		// user in a cookie.
		"cid": {uuid.Must(uuid.NewV4()).String()},
		"t":   {"event"},
		"ec":  {category},
		"ea":  {action},
		"ua":  {r.UserAgent()},
	}
	if label != "" {
		v.Set("el", label)
	}
	if value != nil {
		v.Set("ev", fmt.Sprintf("%d", *value))
	}
	if remoteIP, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		v.Set("uip", remoteIP)
	}
	// NOTE: Google Analytics returns a 200, even if the request is malformed.
        log.Println("%v\n", v);
	_, err := http.PostForm("https://www.google-analytics.com/collect", v)
	return err
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
