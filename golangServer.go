package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var passwordMap map[int]string

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// handler wrapper to collect statistics
func statsCollector(handler http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, req *http.Request) {
		log.Println("Before")
		handler.ServeHTTP(responseWriter, req) // call original
		log.Println("After")
	}
}

type App struct {

	// Top level handlers for the app
	hashHandler *HashHandler
	statsHandler *StatsHandler
	shutdownHandler *ShutdownHandler

}

func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "hash":
		h.hashHandler.ServeHTTP(res, req)
		return
	case "stats":
		h.statsHandler.ServeHTTP(res, req)
	case "shutdown":
		h.shutdownHandler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}

type HashHandler struct {
	// Handlers under the /hash path go here
}

func (h *HashHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		switch req.Method {
		case "GET":
			http.ServeFile(res, req, "hash.html")
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := req.ParseForm(); err != nil {
				fmt.Fprintf(res, "ParseForm() err: %v", err)
				return
			}
			password := req.FormValue("password")
			fmt.Fprintf(res, "Request number goes here... but in the mean time, the password I recieved is %s\n", password)
		default:
			http.Error(res, "405 method not allowed.", http.StatusMethodNotAllowed)
		}
	} else {
		var head string
		head, req.URL.Path = ShiftPath(req.URL.Path)
		hashID, err := strconv.Atoi(head)
		if err != nil {
			http.Error(res, fmt.Sprintf("Invalid hash id %q", head), http.StatusBadRequest)
			return
		}
		switch req.Method {
		case "GET":
			fmt.Fprintf(res, "requested password: %d\n", hashID)
		//	h.handleGet(id)
		//case "PUT":
		//	h.handlePut(id)
		default:
			http.Error(res, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

type StatsHandler struct {
	// Handlers under the /stats path go here
}

func (h *StatsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "stats")
}

type ShutdownHandler struct {
	// Handlers under the /shutdown path go here
}

func (h *ShutdownHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "shutdown")
}

// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
//func healthCheckHandler(responseWriter http.ResponseWriter, req *http.Request) {
//	// A very simple health check.
//	responseWriter.WriteHeader(http.StatusOK)
//	responseWriter.Header().Set("Content-Type", "application/json")
//
//	// In the future we could report back on the status of our DB, or our cache
//	// (e.g. Redis) by performing a simple PING, and include them in the response.
//	fmt.Fprintf(responseWriter, `{"alive": true}`)
//}

func hashAndInsert(password string, reqNumber int, delaySeconds int) {
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	passwordHash := hash(password)
	passwordMap[reqNumber] = passwordHash
}

// the hash function hashes a string with sha512.sum512 and returns it as a base64 string
func hash(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	b := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

// Main function for golang_server
func main() {

	//use the port passed in if it is valid, otherwise use 8080
	port := 8080
	if len(os.Args) > 1 {
		portArg, err := strconv.Atoi(os.Args[1])
		if err == nil {
			if portArg == 80 || (portArg >= 1024 && portArg <= 65535) {
				port = portArg
			}
		}
	}

	passwordMap = make(map[int]string)

	a := &App{
		hashHandler: new(HashHandler),
		statsHandler: new(StatsHandler),
		shutdownHandler: new(ShutdownHandler),
	}

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), a))
}