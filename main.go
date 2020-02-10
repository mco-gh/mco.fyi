package main

import (
  "log"

  firebase "firebase.google.com/go"
  "google.golang.org/api/option"
)

// Use the application default credentials
ctx := context.Background()
conf := &firebase.Config{ProjectID: projectID}
app, err := firebase.NewApp(ctx, conf)
if err != nil {
  log.Fatalln(err)
}

client, err := app.Firestore(ctx)
if err != nil {
  log.Fatalln(err)
}
defer client.Close()

doc := client.Collection("users").Documents(ctx).get();
for {
        doc, err := iter.Next()
        if err == iterator.Done {
                break
        }
        if err != nil {
                log.Fatalf("Failed to iterate: %v", err)
        }
        fmt.Println(doc.Data())
}

import (
    "log"
    "net/http"
)

func redirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "http://www.google.com", 301)
}

func main() {
    http.HandleFunc("/", redirect)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
