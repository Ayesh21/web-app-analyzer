package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
	"web-app-analyzer/internal/logging"
	"web-app-analyzer/web/templates"

	"web-app-analyzer/internal/service"
)

// Load the HTML template files
var templatesHTML = template.Must(template.ParseFS(templates.FS, "index.html", "results.html"))

// HomePageHandler renders the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("Rendering home page")

	// Retrieve error message from query parameters
	errorMessage := r.URL.Query().Get("error")

	// Pass error message to the template
	templatesHTML.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"ErrorMessage": errorMessage,
	})
}

// AnalyzerHandler processes the submitted URLs and validate them
func AnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logging.Logger.Info("Received analysis request")

	// Ensure only GET requests are allowed
	if r.Method != http.MethodGet {
		logging.Logger.Warn("Invalid request method", "method", r.Method)
		http.Redirect(w, r, "/?error=Invalid+Request+Method", http.StatusSeeOther)
		return
	}

	// Retrieve URL from query parameters
	urlValue := r.URL.Query().Get("url")
	fmt.Println("Processing URL:", urlValue)

	if urlValue == "" {
		logging.Logger.Warn("User did not provide a URL")
		http.Redirect(w, r, "/?error=Please+Enter+a+URL", http.StatusSeeOther)
		return
	}

	logging.Logger.Info("Validating URL format", "url", urlValue)

	// Validate the URL format
	parsedURL, err := url.Parse(urlValue)
	if err != nil || !parsedURL.IsAbs() {
		logging.Logger.Warn("Invalid URL format", "url", urlValue)
		http.Redirect(w, r, "/?error=Invalid+URL+Format", http.StatusSeeOther)
		return
	}

	logging.Logger.Info("Fetching URL", "url", parsedURL.String())

	// Fetch the web page content
	response, err := http.Get(parsedURL.String())
	if err != nil {
		logging.Logger.Error("Failed to fetch URL", "url", parsedURL.String(), "error", err)
		http.Redirect(w, r, "/?error=Failed+to+fetch+URL", http.StatusSeeOther)
		return
	}
	defer response.Body.Close()

	// Handle non-200 HTTP responses
	if response.StatusCode != http.StatusOK {
		logging.Logger.Warn("HTTP Error", "status_code", response.StatusCode, "url", parsedURL.String())
		http.Redirect(w, r, fmt.Sprintf("/?error=HTTP+Error+%d", response.StatusCode), http.StatusSeeOther)
		return
	}

	//renders the service method
	pageData := service.AnalyzeHTML(response, parsedURL)
	pageData.URL = urlValue

	// Render results page
	templatesHTML.ExecuteTemplate(w, "results.html", pageData)
	logging.Logger.Info("Analysis completed", "url", urlValue, "time_taken", time.Since(startTime))
}
