package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
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
	http.HandleFunc("/", Home)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

// Home - Home http request
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	log.Println("Get Home")
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	var response string
	for shortLink, link := range linkList {
		response += fmt.Sprintf("Link: <a href=\"http://localhost:9000/short/%s\">http://localhost:9000/short/%s</a> \t\t ShortLink: %s", shortLink, shortLink, link)
	}
	fmt.Fprintf(w, "<h2>Hello and Welcome to the Go URL Shortener!<h2><br>\n")
	fmt.Fprintf(w, response)
	return
}

// addLink - Add a link to the linkList and generate a shorter link
// Example: localhost:9000/addLink?link=https://google.com
func addLink(w http.ResponseWriter, r *http.Request) {
	log.Println("Add Link")
	key, ok := r.URL.Query()["link"]
	if ok {
		if !validLink(key[0]) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not create shortlink need absolute path link. Ex: /addLink?link=https://github.com/")
			return
		}
		log.Println(key)
		if _, ok := linkList[key[0]]; !ok {
			genString := randStringBytes(10)
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

// validLink - check that the link we're creating a shortlink for is a absolute URL path
func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	log.Printf("Checking for valid link: %s", link)
	// Check if string matches the regex
	if r.MatchString(link) {
		return true
	}
	return false
}

// getLink - Find link that matches the shortened link in the linkList
func getLink(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	log.Println("Get Link:", path)
	pathArgs := strings.Split(path, "/")
	if len(pathArgs[2]) < 1 {
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "http://localhost:9000/", http.StatusTemporaryRedirect)
		return
	}
	log.Printf("Redirected to: %s", linkList[pathArgs[2]])
	http.Redirect(w, r, linkList[pathArgs[2]], http.StatusTemporaryRedirect)
	return
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randStringBytes - Create random short link
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
