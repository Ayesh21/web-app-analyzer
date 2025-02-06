package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"web-app-analyzer/internal/model"
	"web-app-analyzer/internal/service"
)

// Load HTML Templates
var templatehtml = template.Must(template.ParseFiles("templates/index.html"))

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	templatehtml.Execute(w, nil)
}

func AnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	urlValue := r.FormValue("url")
	if urlValue == "" {
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Please Enter a URl"})
		return
	}

	parsedURL, err := url.Parse(urlValue)
	if err != nil || !parsedURL.IsAbs() {
		templatehtml.Execute(w, model.PageData{ErrorMessage: "Invalid URL Format"})
		return
	}

	response, err := http.Get(parsedURL.String())
	if err != nil {
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("Failed to fetch URL: %v", err)})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		templatehtml.Execute(w, model.PageData{ErrorMessage: fmt.Sprintf("HTTP Error: %d %s", response.StatusCode, http.StatusText(response.StatusCode))})
		return
	}

	pageData := service.AnalyzeHTML(response, parsedURL)
	pageData.URL = urlValue

	templatehtml.Execute(w, pageData)
}
