package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"web-app-analyzer/internal/model"
	"web-app-analyzer/internal/service"
	"web-app-analyzer/internal/templates"
)

// Load the HTML template file and store it for rendering responses
var templatehtml = template.Must(template.ParseFS(templates.FS, "index.html"))

// HomePageHandler renders the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	templatehtml.Execute(w, nil)
}

// AnalyzerHandler processes user-submitted URLs and performs analysis
func AnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve the URL from the request form
	urlValue := r.FormValue("url")
	if urlValue == "" {
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Please Enter a URl"})
		return
	}

	// Validate the URL format
	parsedURL, err := url.Parse(urlValue)
	if err != nil || !parsedURL.IsAbs() {
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Invalid URL Format"})
		return
	}

	// Fetch the web page content
	response, err := http.Get(parsedURL.String())
	if err != nil {
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("Failed to fetch URL: %v", err)})
		return
	}
	defer response.Body.Close()

	// Handle non-200 HTTP responses
	if response.StatusCode != http.StatusOK {
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("HTTP Error: %d %s", response.StatusCode, http.StatusText(response.StatusCode))})
		return
	}

	//renders the service method
	pageData := service.AnalyzeHTML(response, parsedURL)
	pageData.URL = urlValue

	templatehtml.Execute(w, pageData)
}
