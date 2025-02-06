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

	var inTitle, inLoginForm bool
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return pageData
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

						if !strings.HasPrefix(link, "http") {
							// Resolve the link to an absolute URL
							linkURL, err := baseURL.Parse(link)
							if err == nil {
								link = linkURL.String()
							}
						}

						if strings.HasPrefix(link, "http") {
							pageData.ExternalLinks++
						} else {
							pageData.InternalLinks++
						}
					}
				}
			case "input":
				for _, attr := range token.Attr {
					if attr.Key == "type" && attr.Val == "password" {
						inLoginForm = true
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

	pageData.HasLoginForm = inLoginForm
	pageData.HTMLVersion = "HTML5"

	return pageData
}
