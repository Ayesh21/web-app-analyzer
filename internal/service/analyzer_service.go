package service

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"web-app-analyzer/internal/model"
)

func AnalyzeHTML(response *http.Response, baseURL *url.URL) model.PageData {
	tokenizer := html.NewTokenizer(response.Body)
	defer response.Body.Close()

	pageData := model.PageData{
		HeadingsCount: make(map[string]int),
	}

	var inTitle bool
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return pageData
		case html.DoctypeToken:
			token := tokenizer.Token()
			pageData.HTMLVersion = DetectHTMLVersion(token)
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				inTitle = true
			case "h1", "h2", "h3", "h4", "h5", "h6":
				pageData.HeadingsCount[token.Data]++
			case "a":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val

						// Check for internal links
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
				for _, attr := range token.Attr {
					if attr.Key == "type" && attr.Val == "password" {
						pageData.HasLoginForm = true
					}
				}

			}
		case html.TextToken:
			if inTitle {
				pageData.Title = tokenizer.Token().Data
				inTitle = false
			}
		}

	}

	return pageData
}

func DetectHTMLVersion(token html.Token) string {
	if token.Type != html.DoctypeToken {
		return "Unknown"
	}

	data := strings.ToLower(token.Data)

	if data == "html" && len(token.Attr) == 0 { // Correct HTML5 check
		return "HTML5"
	}

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
