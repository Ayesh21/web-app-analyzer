package service

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"web-app-analyzer/internal/model"
)

// Mock HTTP response to simulate HTML content.
func mockResponse(htmlContent string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(htmlContent)),
	}
}

func TestAnalyzeHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		baseURL  string
		expected model.PageData
	}{
		{
			name:    "Test with internal and external links",
			html:    `<html><head><title>Test Page</title></head><body><a href="https://test.com">External Link</a><a href="/internal">Internal Link</a></body></html>`,
			baseURL: "http://localhost",
			expected: model.PageData{
				Title:         "Test Page",
				HeadingsCount: map[string]int{},
				ExternalLinks: 1,
				InternalLinks: 1,
				HasLoginForm:  false,
				HTMLVersion:   "",
			},
		},
		{
			name:    "Test with login form",
			html:    `<html><body><form><input type="password" /></form></body></html>`,
			baseURL: "http://localhost",
			expected: model.PageData{
				Title:         "",
				HeadingsCount: map[string]int{},
				ExternalLinks: 0,
				InternalLinks: 0,
				HasLoginForm:  true,
				HTMLVersion:   "",
			},
		},
		{
			name:    "Test with headings",
			html:    `<html><body><h1>Heading 1</h1><h2>Heading 2</h2></body></html>`,
			baseURL: "http://localhost",
			expected: model.PageData{
				Title:         "",
				HeadingsCount: map[string]int{"h1": 1, "h2": 1},
				ExternalLinks: 0,
				InternalLinks: 0,
				HasLoginForm:  false,
				HTMLVersion:   "",
			},
		},
		{
			name:    "Test HTML5 DOCTYPE declaration",
			html:    `<!DOCTYPE html><html><head><title>HTML5 Page</title></head><body></body></html>`,
			baseURL: "http://localhost",
			expected: model.PageData{
				Title:         "HTML5 Page",
				HeadingsCount: map[string]int{},
				ExternalLinks: 0,
				InternalLinks: 0,
				HasLoginForm:  false,
				HTMLVersion:   "HTML5", // Expecting "HTML5" due to the DOCTYPE
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert the base URL string to a *url.URL object
			baseURL, err := url.Parse(tt.baseURL)
			assert.NoError(t, err)

			// Create the mock HTTP response
			response := mockResponse(tt.html)

			// Call the AnalyzeHTML function
			actual := AnalyzeHTML(response, baseURL)

			// Assert that the actual result matches the expected
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		doctype  string
		expected string
	}{
		{
			name:     "HTML5 DOCTYPE",
			doctype:  "<!DOCTYPE html>",
			expected: "HTML5",
		},
		{
			name:     "XHTML 1.0 Strict DOCTYPE",
			doctype:  `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
			expected: "XHTML",
		},
		{
			name:     "HTML 4.01 Strict DOCTYPE",
			doctype:  `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">`,
			expected: "HTML 4.01",
		},
		{
			name:     "HTML 4.01 Transitional DOCTYPE",
			doctype:  `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">`,
			expected: "HTML 4.01",
		},
		{
			name:     "HTML 4.01 Frameset DOCTYPE",
			doctype:  `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Frameset//EN" "http://www.w3.org/TR/html4/frameset.dtd">`,
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML 1.0 Transitional DOCTYPE",
			doctype:  `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
			expected: "XHTML",
		},
		{
			name:     "XHTML 1.0 Frameset DOCTYPE",
			doctype:  `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`,
			expected: "XHTML",
		},
		{
			name:     "Unknown DOCTYPE",
			doctype:  `<!DOCTYPE unknown>`,
			expected: "Unknown",
		},
		{
			name:     "Empty DOCTYPE",
			doctype:  `<!DOCTYPE>`,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token html.Token

			if tt.expected == "HTML5" {
				token = html.Token{
					Type: html.DoctypeToken,
					Data: "html", // Crucial for HTML5
				}
			} else { // Handle other doctypes which have public/system identifiers
				token = html.Token{
					Type: html.DoctypeToken,
					Data: "html",
					Attr: []html.Attribute{
						{Key: "public", Val: tt.doctype},
					},
				}
			}

			actual := DetectHTMLVersion(token)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
