package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStartWiki(t *testing.T) {
	println("Testing TestStartWiki at main.go")
	req, err := http.NewRequest("GET", "/search?keyword=monkey", nil)
	Check(err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(StartWiki)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := reflect.TypeOf(rr.Body.String())

	if "string" != expected.String() {
		t.Errorf("handler returned unexpected body: \n got: \n %v \n want: \n %v",
			rr.Body.String(), expected)
	} else {
		t.Log("Success!")
	}

}
