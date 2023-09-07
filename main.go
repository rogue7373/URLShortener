package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	linkList map[string]string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	linkList = map[string]string{}

	http.HandleFunc("/addLink", addLink)
	http.HandleFunc("/short/", getLink)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

// addLink - Add a link to the linkList and generate a shorter link
func addLink(w http.ResponseWriter, r *http.Request) {
	key, ok := r.URL.Query()["link"]
	if ok {
		if _, ok := linkList[key[0]]; !ok {
			genString := fmt.Sprint(rand.Int63n(1000))
			linkList[genString] = key[0]
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusAccepted)
			linkString := fmt.Sprintf("<a href=\"http://localhost:9000/short/%s\">http://localhost:9000/short/%s</a>", genString, genString)
			fmt.Fprintf(w, "Added shortlink\n")
			fmt.Fprintf(w, linkString)
			return
		}
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Already have this link")
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Failed to add link")
	return
}

// getLink - Find link that matches the shortened link in the linkList
func getLink(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")
	log.Printf("Redirected to: %s", linkList[pathArgs[2]])
	http.Redirect(w, r, linkList[pathArgs[2]], http.StatusPermanentRedirect)
	return
}
