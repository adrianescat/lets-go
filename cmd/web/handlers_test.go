package main

import (
	"bytes"
	"github.com/adrianescat/lets-go/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder.
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request.
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		// Typically you should call t.Fatal() in situations where it doesn't make sense to
		// continue the current test. It will mark the test as failed, log the error, and then
		// completely stop execution of the current test (or sub-test).
		t.Fatal(err)
	}

	// Call the ping handler function, passing in the
	// httptest.ResponseRecorder and http.Request.
	ping(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the
	// http.Response generated by the ping handler.
	rs := rr.Result()

	// Check that the status code written by the ping handler was 200.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	// And we can check that the response body written by the ping handler
	// equals "Pong".
	assert.Equal(t, string(body), "Pong")
}
