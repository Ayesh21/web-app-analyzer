package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHomePageHandler ensures the home page loads successfully.
func TestHomePageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HomePageHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
}

// TestAnalyzerHandler_InvalidMethod ensures GET requests are redirected.
func TestAnalyzerHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/analyze", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
}

// TestAnalyzerHandler_EmptyURL ensures that an empty URL triggers an error.
func TestAnalyzerHandler_EmptyURL(t *testing.T) {
	req, err := http.NewRequest("POST", "/analyze", strings.NewReader("url="))
	assert.NoError(t, err, "Error creating request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
	assert.Contains(t, recorder.Body.String(), "Please Enter a URl", "Expected error message for empty URL")
}

// TestAnalyzerHandler_InvalidURL ensures that an invalid URL triggers an error.
func TestAnalyzerHandler_InvalidURL(t *testing.T) {
	req, err := http.NewRequest("POST", "/analyze", strings.NewReader("url=invalid-url"))
	assert.NoError(t, err, "Error creating request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
	assert.Contains(t, recorder.Body.String(), "Invalid URL Format", "Expected error message for invalid URL")
}

// TestAnalyzerHandler_FailedFetch simulates a request to an unreachable URL.
func TestAnalyzerHandler_FailedFetch(t *testing.T) {
	req, err := http.NewRequest("POST", "/analyze", strings.NewReader("url=http://invalid.test.url"))
	assert.NoError(t, err, "Error creating request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
	assert.Contains(t, recorder.Body.String(), "Failed to fetch URL", "Expected error message for failed fetch")
}
