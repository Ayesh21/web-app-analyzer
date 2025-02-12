package service

import (
	"net/http"
	"net/url"
	"strings"
	"web-app-analyzer/internal/logging"

	"golang.org/x/net/html"
	"web-app-analyzer/internal/model"
)

// AnalyzeHTML parses the HTML content and extracts metadata
func AnalyzeHTML(response *http.Response, baseURL *url.URL) model.PageData {
	logging.Logger.Info("Starting HTML analysis", "url", baseURL.String())

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
			logging.Logger.Info("Reached end of HTML document")
			return pageData
		case html.DoctypeToken:
			token := tokenizer.Token()
			// Extract and detect the HTML version
			pageData.HTMLVersion = DetectHTMLVersion(token)
			logging.Logger.Info("Detected HTML version", "version", pageData.HTMLVersion)
		case html.StartTagToken, html.SelfClosingTagToken:
			// Handle different HTML elements based on tag type
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				// Mark that we are inside a title tag
				inTitle = true
			case "h1", "h2", "h3", "h4", "h5", "h6":
				// Count occurrences of each heading type
				pageData.HeadingsCount[token.Data]++
				logging.Logger.Info("Found heading", "tag", token.Data)
			case "a":
				// Extract link attributes
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
							logging.Logger.Info("Found internal link", "url", link)
						} else {
							pageData.ExternalLinks++
							logging.Logger.Info("Found external link", "url", link)
						}
					}
				}
			case "input":
				// Check if the page has a login form by detecting password fields
				for _, attr := range token.Attr {
					if attr.Key == "type" && attr.Val == "password" {
						pageData.HasLoginForm = true
						logging.Logger.Info("Detected a login form")
					}
				}

			}
		case html.TextToken:
			// Capture the title text when inside a <title> tag
			if inTitle {
				pageData.Title = tokenizer.Token().Data
				logging.Logger.Info("Extracted page title", "title", pageData.Title)
				inTitle = false
			}
		}

	}

	return pageData
}

// DetectHTMLVersion determines the HTML version based on the DOCTYPE declaration.
func DetectHTMLVersion(token html.Token) string {
	logging.Logger.Info("Processing HTML version")
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

	logging.Logger.Warn("Could not determine HTML version")
	return "Unknown"
}
