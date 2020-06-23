// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"cloud.google.com/go/firestore"
)

var (
	mu           sync.Mutex
	linkData     map[string]interface{}
	renderedHome []byte

	doc *firestore.DocumentRef
)

var homeTemplate = template.Must(template.ParseFiles("home.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	resp := renderedHome
	mu.Unlock()
	w.Write(resp)
}

// renderHome must only be called while mu is held.
func renderHome() ([]byte, error) {
	type link struct {
		Key   string // short name
		Count int64  // count
		URL   string // URL
		Desc  string // description
	}
	var links []link

	for short, v := range linkData {
		v, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		private := v["private"]
		if private != nil && private.(bool) == true {
			continue
		}
		count, _ := v["count"].(int64)
		desc, _ := v["desc"].(string)
		links = append(links, link{
			Key:   short,
			Count: count,
			URL:   v["url"].(string),
			Desc:  desc,
		})
	}
	sort.Slice(links, func(i, j int) bool {
		return links[i].Count > links[j].Count
	})

	var buf bytes.Buffer
	if err := homeTemplate.Execute(&buf, links); err != nil {
		log.Printf("Template.Execute: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")

	mu.Lock()
	link, found := linkData[path]
	mu.Unlock()

	if !found {
		http.ServeFile(w, r, "404.html")
		return
	}

	linkDoc := link.(map[string]interface{})
	u, ok := linkDoc["url"]
	if !ok {
		log.Println(w, "no URL found for event: %v", path)
		// TODO: return a response?
		return
	}

	go func() {
		ctx := context.Background()
		log.Printf("before: %d", linkDoc["count"])
		prevCount, _ := linkDoc["count"].(int64)
		linkDoc["count"] = prevCount + 1
		log.Printf("after: %d", linkDoc["count"])
		// TODO: this races, use firestore.Increment
		doc.Set(ctx, linkData)
	}()

	http.Redirect(w, r, u.(string), 301)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if path == "" || path == "/" {
		homeHandler(w, r)
		return
	}

	if strings.HasPrefix(path, "css/") || strings.HasPrefix(path, "img/") {
		http.ServeFile(w, r, path)
		return
	}

	linkHandler(w, r)
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
	docSnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	linkData = docSnap.Data()
	renderedHome, _ = renderHome()

	// This function runs in the background. It gets notified
	// anytime the dataset changes, and reloads the local copy
	// in response to those notifications so that the running
	// instance always has the latest version of the data handy.
	go func() {
		iter := doc.Snapshots(ctx)
		defer iter.Stop()
		for {
			docSnap, err := iter.Next()
			if err != nil {
				log.Fatalln(err)
			}
			mu.Lock()
			linkData = docSnap.Data()
			renderedHome, _ = renderHome()
			mu.Unlock()
		}
	}()

	http.HandleFunc("/", rootHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
