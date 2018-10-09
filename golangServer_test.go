package main

import (
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

//func TestHealthCheckHandler(t *testing.T) {
//	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
//	// pass 'nil' as the third parameter.
//	req, err := http.NewRequest("GET", "/health-check", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(healthCheckHandler)
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if rr.Code != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			rr.Code, http.StatusOK)
//	}
//
//	// Check the response body is what we expect.
//	expected := `{"alive": true}`
//	if rr.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			rr.Body.String(), expected)
//	}
//}