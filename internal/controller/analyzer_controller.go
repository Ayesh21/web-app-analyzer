package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
	"web-app-analyzer/internal/logging"
	"web-app-analyzer/web/templates"

	"web-app-analyzer/internal/model"
	"web-app-analyzer/internal/service"
)

// Load the HTML template file and store it for rendering responses
var templatehtml = template.Must(template.ParseFS(templates.FS, "index.html"))

// HomePageHandler renders the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("Rendering home page")
	templatehtml.Execute(w, nil)
}

// AnalyzerHandler processes user-submitted URLs and performs analysis
func AnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logging.Logger.Info("Received analysis request")

	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		logging.Logger.Warn("Invalid request method", "method", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve the URL from the request form
	urlValue := r.FormValue("url")
	if urlValue == "" {
		logging.Logger.Warn("User did not provide a URL")
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Please Enter a URl"})
		return
	}

	logging.Logger.Info("Processing URL", "url", urlValue)

	// Validate the URL format
	parsedURL, err := url.Parse(urlValue)
	if err != nil || !parsedURL.IsAbs() {
		logging.Logger.Warn("Invalid URL format", "url", urlValue)
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Invalid URL Format"})
		return
	}

	logging.Logger.Info("Fetching URL", "url", parsedURL.String())

	// Fetch the web page content
	response, err := http.Get(parsedURL.String())
	if err != nil {
		logging.Logger.Error("Failed to fetch URL", "url", parsedURL.String(), "error", err)
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("Failed to fetch URL: %v", err)})
		return
	}
	defer response.Body.Close()

	// Handle non-200 HTTP responses
	if response.StatusCode != http.StatusOK {
		logging.Logger.Warn("HTTP Error", "status_code", response.StatusCode, "url", parsedURL.String())
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("HTTP Error: %d %s", response.StatusCode, http.StatusText(response.StatusCode))})
		return
	}

	//renders the service method
	pageData := service.AnalyzeHTML(response, parsedURL)
	pageData.URL = urlValue

	templatehtml.Execute(w, pageData)
	logging.Logger.Info("Analysis completed", "url", urlValue, "time_taken", time.Since(startTime))

}
