package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Global vars
var passwordMap map[int]string
var postCount Counter
var postStats []int64
var shutdown bool
var hs *http.Server


// Syncronized counter
type Counter struct {
	mu  sync.Mutex
	count   int
}

// Increment and return the current count
func (c *Counter) Increment() (count int) {
	c.mu.Lock()
	c.count += 1
	count = c.count
	c.mu.Unlock()
	return
}

// Get current count
func (c *Counter) Count() (count int) {
	return
}

// Get time in microseconds
// getting time in microseconds, due to the fact that
// milliseconds seem to not provide enough resolution
// for this small server
func makeTimestamp() int64 {
	return time.Now().UnixNano()/int64(time.Microsecond)
}

// ShiftPath splits off the first component of p
// to the head, and all else in the tail,
// if there is only one element in the path, tail
// will contain a "/"
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// The main handler for the application
type App struct {

	// Top level handlers for the app
	hashHandler *HashHandler
	statsHandler *StatsHandler
	shutdownHandler *ShutdownHandler

}

// ServeHTTP method for App
func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "hash":
		before := makeTimestamp()
		h.hashHandler.ServeHTTP(res, req)
		after := makeTimestamp()
		if req.Method == "POST" && req.URL.Path == "/" {
			postStats = append(postStats, (after-before))
			log.Print(postStats)
			log.Printf("before: %d", before)
			log.Printf("after:  %d", after)
		}
	case "stats":
		h.statsHandler.ServeHTTP(res, req)
	case "shutdown":
		h.shutdownHandler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}
// HashHandler struct, Handlers under the '/hash' path go here
// This is also a data structure to attach the ServeHTTP method to
type HashHandler struct {
}

// ServeHTTP method for HashHandler
func (h *HashHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		switch req.Method {
		case "GET":
			http.ServeFile(res, req, "hash.html")
		case "POST":
			if err := req.ParseForm(); err != nil {
				fmt.Fprintf(res, "ParseForm() err: %v", err)
				return
			}
			password := req.FormValue("password")
			currentCount := postCount.Increment()
			fmt.Fprintf(res, strconv.Itoa(currentCount))
			go hashAndInsert(password, currentCount, 5)
		default:
			http.Error(res, "405 method not allowed.", http.StatusMethodNotAllowed)
		}
	} else {
		var head string
		head, req.URL.Path = ShiftPath(req.URL.Path)
		hashID, err := strconv.Atoi(head)
		if err != nil {
			http.Error(res, fmt.Sprintf("Invalid hash id %s", head), http.StatusBadRequest)
			return
		}
		switch req.Method {
		case "GET":
			hash := passwordMap[hashID]
			if "" != hash {
				fmt.Fprintf(res, hash)
			} else {
				http.Error(res, "404 resource not found", http.StatusNotFound)
			}

		default:
			http.Error(res, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

// StatsHandler object/struct
type StatsHandler struct {
	// Handlers under the /stats path go here
}

// ServeHTTP method for StatsHandler
func (h *StatsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	total := postCount.Count()
	var average float64
	var sum int64
	average = 0
	if len(postStats) != 0 {
		for i := 0; i < len(postStats); i++ {
			sum += postStats[i]
		}
		//Response times were well under a millisecond, so this...
		average = float64(sum / int64(total)) / 1000.0
	}

	fmt.Fprintf(res, `{"total": %d, "average": %f}`, total, average)
}

// ShutdownHandler object/struct
type ShutdownHandler struct {
	// Handlers under the /shutdown path go here
}

// ServerHTTP method for ShutdownHandler
func (h *ShutdownHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "shutdown recieved, shutting down")
	go graceful()
}

// Wait for delay, hash the password, and store it in the passwordMap
func hashAndInsert(password string, reqNumber int, delaySeconds int) {
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	passwordHash := hash(password)
	passwordMap[reqNumber] = passwordHash
}

// the hash function hashes a string with sha512.sum512 and returns it as a base64 string
func hash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	b := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}

// Stop the server gracefully, and kill it after a timeout if needed
func graceful() {
	timeout := 10*time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
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

	shutdown = false
	passwordMap = make(map[int]string)
	postCount.count = 0
	postStats = make([]int64, 0)

	app := &App{
		hashHandler: new(HashHandler),
		statsHandler: new(StatsHandler),
		shutdownHandler: new(ShutdownHandler),
	}
	hs = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: app}


	log.Printf("Listening on http://0.0.0.0%s\n", hs.Addr)

	if err := hs.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}