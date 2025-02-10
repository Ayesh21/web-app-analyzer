package main

import (
	"fmt"
	"log"
	"net/http"

	"web-app-analyzer/internal/controller"
)

func main() {
	// Route for home page
	http.HandleFunc("/", controller.HomePageHandler)

	// Route for analyzing web pages
	http.HandleFunc("/analyze", controller.AnalyzerHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
