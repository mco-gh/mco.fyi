package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

var linkdata map[string]interface{}
var doc *firestore.DocumentRef
var gaPropertyID = "UA-158788691-1"

type kv struct {
	K string
	C int64
	U string
	D string
}

func redirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimLeft(r.URL.Path, "/")
	if path == "" || path == "/" {
		t, err := template.ParseFiles("home.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
		var kvs []kv
		for k, v := range linkdata {
			tmp := v.(map[string]interface{})
			count := tmp["count"].(int64)
			desturl := tmp["url"].(string)
			desc := tmp["desc"].(string)
			kvs = append(kvs, kv{k, count, desturl, desc})
		}
		sort.Slice(kvs, func(i, j int) bool {
			return kvs[i].C > kvs[j].C
		})
		err = t.Execute(w, kvs)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else if strings.HasPrefix(path, "css/") ||
		strings.HasPrefix(path, "img/") {
		http.ServeFile(w, r, path)
	} else if m, ok := linkdata[path]; ok {
		v := m.(map[string]interface{})
		if u, ok := v["url"]; ok {
			log.Println("before: %d", v["count"])
			v["count"] = v["count"].(int64) + 1
			log.Println("after: %d", v["count"])
			doc.Set(ctx, linkdata)
			http.Redirect(w, r, u.(string), 301)
		} else {
			log.Println(w, "no URL found for event: %v", path)
			return
		}
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
	defer client.Close()
	doc = client.Doc("Redirects/Shortlinks")
	docsnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//updateTime := docsnap.UpdateTime
	linkdata = docsnap.Data()

	//iter := shortlinks.Snapshots(ctx)
	//defer iter.Stop()
	//for {
	//docsnap, err := iter.Next()
	//if err != nil {
	// TODO: Handle error.
	//}
	//_ = docsnap // TODO: Use DocumentSnapshot.
	//}

	http.HandleFunc("/", redirect)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
