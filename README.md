# Web-Page-Analyzer
## Overview

Web App Analyzer is a Go-based web application that analyzes and processes URLs provided by users. It fetches web pages, extracts key data, and displays results in an HTML-based user interface.

Features

* Analyze web pages by submitting URLs.
* Validate URL input and handle errors gracefully.
* Serve responses using pre-defined HTML templates.
* Unit-tested on controller, service layers and test coverage

structure 

## Project Structure

/web-app-analyzer  
│── go.mod                  # Go module dependencies  
│── go.sum                  # Go dependencies  
│── main.go                 # Application entry point  
│── /.github  
│   ├── /workflows  
│   │   ├── development.yml   # GitHub action CI/CD pipeline configurations  
│── /internal  
│   │── /controller  
│   │   ├── analyzer_controller.go          # Controller logic  
│   │   ├── analyzer_controller_test.go     # Unit tests for controller  
│   │── /service  
│   │   ├── analyzer_service.go             # Business logic layer  
│   │   ├── analyzer_service_test.go        # Unit tests for service  
│   │── /model  
│   │   ├── page_data.go                    # Data layer  
│   │── /templates  
│   │   ├── index.html                      # HTML template for rendering responses  
│── DockerFile                 # Docker file  
│── README.md                  # Project overview  

## Installation
Ensure you have Go 1.23 installed.
```
git clone https://github.com/Ayesh21/web-app-analyzer.git
cd web-app-analyzer
go mod tidy
```
## Building & Running the Application
`go build`
`go run .`

## Running Tests
`go test -v ./internal/controller/... ./internal/service/...`

## Generating Test Coverage Report
```
go test ./internal/controller/... ./internal/service/... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
start coverage.html
```
## API Endpoints
`GET /`

* Serves the home page.
* `http://localhost:8080/`

`POST /analyze`

* Accepts a URL for analysis.
* Returns extracted page data or an error message.










