package model

// PageData /** PageData holds the analysis results for a web page,It includes metadata such as URL,
// HTML version, title, heading counts, link counts, login form presence, and errors.
type PageData struct {
	URL           string
	HTMLVersion   string
	Title         string
	HeadingsCount map[string]int
	InternalLinks int
	ExternalLinks int
	HasLoginForm  bool
	ErrorMessage  string
}
