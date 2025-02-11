package service

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"web-app-analyzer/internal/model"
)

// AnalyzeHTML parses the HTML content and extracts metadata
func AnalyzeHTML(response *http.Response, baseURL *url.URL) model.PageData {
	// Initialize the HTML tokenizer to parse the response body
	tokenizer := html.NewTokenizer(response.Body)
	defer response.Body.Close()

	// Initialize a PageData object to store extracted information
	pageData := model.PageData{
		HeadingsCount: make(map[string]int),
	}

	// Flag to track if we are inside a <title> tag
	var inTitle bool

	// Iterate over HTML tokens
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// Reached end of document; return the collected data
			log.Println("Reached end of HTML document")
			return pageData
		case html.DoctypeToken:
			token := tokenizer.Token()
			// Extract and detect the HTML version
			log.Printf("Extract and detect HTML version: %s\n", pageData.HTMLVersion)
			pageData.HTMLVersion = DetectHTMLVersion(token)
		case html.StartTagToken, html.SelfClosingTagToken:
			// Handle different HTML elements based on tag type
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				// Mark that we are inside a title tag
				inTitle = true
			case "h1", "h2", "h3", "h4", "h5", "h6":
				// Count occurrences of each heading type
				log.Printf("Found heading: %s\n", token.Data)
				pageData.HeadingsCount[token.Data]++
			case "a":
				// Extract link attributes
				log.Println("Processing external and internal links")
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val

						// Check if the link is internal or external
						if !strings.HasPrefix(link, "http") && !strings.HasPrefix(link, "//") {
							linkURL, err := baseURL.Parse(link)
							if err == nil {
								link = linkURL.String()
							}
							pageData.InternalLinks++
						} else {
							pageData.ExternalLinks++
						}
					}
				}
			case "input":
				// Check if the page has a login form by detecting password fields
				for _, attr := range token.Attr {
					if attr.Key == "type" && attr.Val == "password" {
						pageData.HasLoginForm = true
						log.Println("Detected a login form")
					}
				}

			}
		case html.TextToken:
			// Capture the title text when inside a <title> tag
			if inTitle {
				pageData.Title = tokenizer.Token().Data
				log.Printf("Extracted page title: %s\n", pageData.Title)
				inTitle = false
			}
		}

	}

	return pageData
}

// DetectHTMLVersion determines the HTML version based on the DOCTYPE declaration.
func DetectHTMLVersion(token html.Token) string {
	log.Println("Processing HTML version")
	if token.Type != html.DoctypeToken {
		return "Unknown"
	}

	data := strings.ToLower(token.Data)

	// If doctype is "html" with no attributes, assume HTML5
	if data == "html" && len(token.Attr) == 0 {
		return "HTML5"
	}

	// Check for other HTML versions based on attributes
	for _, attr := range token.Attr {
		val := strings.ToLower(attr.Val)
		if strings.Contains(val, "xhtml") {
			return "XHTML"
		} else if strings.Contains(val, "html 4.01") {
			return "HTML 4.01"
		}
	}
	return "Unknown"
}
