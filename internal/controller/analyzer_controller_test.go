package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"web-app-analyzer/internal/logging"

	"github.com/stretchr/testify/assert"
)

// TestMain runs before all tests
func TestMain(m *testing.M) {
	logging.InitLogger()
	m.Run()
}

// TestHomePageHandler Testing whether the home page loads successfully.
func TestHomePageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HomePageHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
}

// TestAnalyzerHandler_InvalidMethod Testing whether non-GET HTTP requests are redirected.
func TestAnalyzerHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("POST", "/results", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
	assert.Equal(t, "/?error=Invalid+Request+Method", recorder.Header().Get("Location"), "Expected redirect to home page with error")
}

// TestAnalyzerHandler_EmptyURL Testing whether that an empty URL triggers a redirect.
func TestAnalyzerHandler_EmptyURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/results", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
	assert.Equal(t, "/?error=Please+Enter+a+URL", recorder.Header().Get("Location"), "Expected redirect with error message")
}

// TestAnalyzerHandler_InvalidURL Testing whether that an invalid URL triggers a redirect.
func TestAnalyzerHandler_InvalidURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/results?url=invalid-url", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
	assert.Equal(t, "/?error=Invalid+URL+Format", recorder.Header().Get("Location"), "Expected redirect with invalid URL error")
}

// TestAnalyzerHandler_FailedFetch simulates a request to an unreachable URL.
func TestAnalyzerHandler_FailedFetch(t *testing.T) {
	req, err := http.NewRequest("GET", "/results?url=http://invalid.test.url", nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
	assert.Equal(t, "/?error=Failed+to+fetch+URL", recorder.Header().Get("Location"), "Expected redirect with failed fetch error")
}

// TestAnalyzerHandler_HttpError simulates a scenario where the requested URL returns a 500 error.
func TestAnalyzerHandler_HttpError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	req, err := http.NewRequest("GET", "/results?url="+url.QueryEscape(mockServer.URL), nil)
	assert.NoError(t, err, "Error creating request")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzerHandler)
	handler.ServeHTTP(recorder, req)

	expectedRedirect := fmt.Sprintf("/?error=HTTP+Error+%d", http.StatusInternalServerError)

	assert.Equal(t, http.StatusSeeOther, recorder.Code, "Expected HTTP status 303 See Other")
	assert.Equal(t, expectedRedirect, recorder.Header().Get("Location"), "Expected redirect with HTTP error message")
}
