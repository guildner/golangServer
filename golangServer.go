package main

import (
	"fmt"
	"log"
	"net/http"
)

// handler for the /hash endpoint
func hashHandler(responseWriter http.ResponseWriter, req *http.Request) {
	passNum := req.URL.Path[len("/hash/")]

	switch req.Method {
	case "GET":
		http.ServeFile(responseWriter, req, "hash.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(responseWriter, "ParseForm() err: %v", err)
			return
		}
		password := req.FormValue("password")
		fmt.Fprintf(responseWriter, "Request number goes here... but in the mean time, the password I recieved is %s\n", password)
	default:
		http.Error(responseWriter, "404 not found.", http.StatusNotFound)
	}
}

// handler for the /stats endpoint

// handler for the /shutdown endpoint

// the hash function

// Main function for golang_server
func main() {
	http.HandleFunc("/hash", hashHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}