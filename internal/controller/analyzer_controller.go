package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"web-app-analyzer/internal/model"
	"web-app-analyzer/internal/service"
	"web-app-analyzer/internal/templates"
)

// Load the HTML template file and store it for rendering responses
var templatehtml = template.Must(template.ParseFS(templates.FS, "index.html"))

// HomePageHandler renders the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering home page")
	templatehtml.Execute(w, nil)
}

// AnalyzerHandler processes user-submitted URLs and performs analysis
func AnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Println("Received request for analysis")

	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		log.Println("Invalid request method, redirecting to home page")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve the URL from the request form
	urlValue := r.FormValue("url")
	if urlValue == "" {
		log.Println("User did not provide a URL")
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Please Enter a URl"})
		return
	}

	// Validate the URL format
	parsedURL, err := url.Parse(urlValue)
	if err != nil || !parsedURL.IsAbs() {
		log.Printf("Invalid URL format: %s\n", urlValue)
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Invalid URL Format"})
		return
	}

	log.Printf("Fetching URL: %s\n", parsedURL.String())
	// Fetch the web page content
	response, err := http.Get(parsedURL.String())
	if err != nil {
		log.Printf("Failed to fetch URL: %v\n", err)
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("Failed to fetch URL: %v", err)})
		return
	}
	defer response.Body.Close()

	// Handle non-200 HTTP responses
	if response.StatusCode != http.StatusOK {
		log.Printf("HTTP Error: %d %s\n", response.StatusCode, http.StatusText(response.StatusCode))
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("HTTP Error: %d %s", response.StatusCode, http.StatusText(response.StatusCode))})
		return
	}

	//renders the service method
	pageData := service.AnalyzeHTML(response, parsedURL)
	pageData.URL = urlValue

	templatehtml.Execute(w, pageData)
	log.Printf("Analysis completed in %v\n", time.Since(startTime))
}
