package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
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

	//TODO: test returned JSON
}