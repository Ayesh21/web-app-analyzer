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
│   │── /logging  
│   │   ├── logger.go                      # Log configuration      
├── /web  
│   ├── /templates  
│   │   ├── index.html                       # HTML template for main UI  
│   │   ├── results.html                       # HTML template for rendering results   
│   ├── /static   
│   │   ├── /css   
│   │   │   ├── styles.css                     # CSS styling   
│   ├── /images   
│   │   ├── background.jpg                     #Background image 
│── /cmd  
│   │── /app  
│   │   ├── main.go                 # Application entry point  
│── /logs  
│   ├── app.log                             # Log file  
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
`go build -o web-app-analyzer ./cmd/app/main.go`

`go run cmd/app/main.go`

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

`GET /analyze`

* Accepts a URL for analysis.
* Returns extracted page data or an error message.

## Errors

<img src="https://github.com/user-attachments/assets/775e5b94-641b-4dbf-9cdf-75f661135e7a" width="300">

<img src="https://github.com/user-attachments/assets/0ab90b88-e22f-42fd-888c-a4ed74d36c7f" width="300">

<img src="https://github.com/user-attachments/assets/7050a68b-7bdf-4612-ad37-60a97e960538" width="300">

## Demo

[Watch the video](https://drive.google.com/file/d/1Y5N-hTQf2ZIaf5a98t19FfePcg-5U_eh/view?usp=sharing)

## Problems Faced During The Implementation 

### Routing Issue
When I ran the unit tests, I got the following error:  `panic: open C:\Users\Ayesh Navindu\OneDrive\Desktop\Go\web-app-analyzer\web-app-analyzer\internal\templates\index.html: The system cannot find the path specified.` However, the application works fine when I run it.

As the solution I used  `//go:embed` (we can use this to embed one or more files into a file system:)

```
//go:embed index.html results.html
var FS embed.FS
```
### Divide Index.html Into Two

To make the final output more readable, I decided to handle a separate `HTML` file for the results. So, I created two `HTML` files: `index.html` and `results.html`.

When I had only `index.html`, the form used `POST` because the same page handled both the input and the result rendering. However, after introducing two `HTML` files, I had to use `GET` since I needed a way to pass data from `index.html` to `results.html`. This is why I switched to `GET`.

After switching to `HTTP GET`, I refactored all the controller and contoller unit tests to retrieve data from query parameters.













