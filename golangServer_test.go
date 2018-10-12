package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestTest(t *testing.T) {
	result := true

	if result == false {
		t.Errorf("Testing a failed test")
	}
}

func TestShiftPath(t *testing.T) {
	testString := "/test/path/for/shift/path/"

	expectedHead := "test"
	expectedTail := "/path/for/shift/path"

	head, tail := ShiftPath(testString)

	if head != expectedHead {
		t.Errorf("Expected %s but got %s", expectedHead, head)
	}
	if tail != expectedTail {
		t.Errorf( "Expected %s but got %s", expectedTail, tail)
	}
}

func TestHash(t *testing.T) {
	password := "angryMonkey"
	expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	hash := hash(password)

	if hash != expectedHash {
		t.Errorf("Expected %s but got %s", expectedHash, hash)
	}
}

func TestHashAndInsert(t *testing.T) {
	passwordMap = make(map[int]string)
	password := "angryMonkey"
	expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	hashAndInsert(password, 1, 0)

	hash := passwordMap[1]

	if hash != expectedHash {
		t.Errorf("Expected %s but got %s", expectedHash, hash)
	}

}

func TestHashHandler(t *testing.T) {
	// Create a request to pass to the handler.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(HashHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body
	buf, err := ioutil.ReadFile("hash.html")
	expected := string(buf)
	actual := rr.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}

func TestHashHandlerPostAndGetWithID(t *testing.T) {
	form := url.Values{}
	form.Add("password", "angryMonkey")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Form = form

	rr := httptest.NewRecorder()
	handler := new(HashHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check response body
	expected := "1"
	actual := rr.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}

	rr = httptest.NewRecorder()

	// Make another request
	handler.ServeHTTP(rr, req)

	// Check response body again
	expected = "2"
	actual = rr.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}

	// Request the hash with ID=1
	req, err = http.NewRequest("GET", "/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the password has has not been stored yet
	expected = "404 resource not found\n"
	actual = rr.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}

	// Try again after sleeping a while
	time.Sleep(6*time.Second)
	req, err = http.NewRequest("GET", "/2", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the first request was hashed and stored
	expected = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	actual = rr.Body.String()
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

}

func TestStatsHandler(t *testing.T) {
	// Create a request to pass to the handler.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(StatsHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

func TestShutdownHandler(t *testing.T) {
	// Create a request to pass to the handler.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(ShutdownHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}
}
