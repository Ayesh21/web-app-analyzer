package main

import (
	"fmt"
	"log"
	"net/http"
	"web-app-analyzer/internal/controller"
	"web-app-analyzer/internal/logging"
)

func main() {
	// Initialize logger
	logging.InitLogger()
	logging.Logger.Info("Starting server...")

	// Route to the static files (CSS, JS, Images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("web/util/images"))))

	// Route for home page
	http.HandleFunc("/", controller.HomePageHandler)

	// Route for analyzing web pages (results page)
	http.HandleFunc("/results", controller.AnalyzerHandler)

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
