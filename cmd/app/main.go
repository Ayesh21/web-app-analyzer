package main

import (
	"fmt"
	"log"
	"net/http"

	"web-app-analyzer/internal/controller"
)

func main() {
	http.HandleFunc("/", controller.HomePageHandler)
	http.HandleFunc("/analyze", controller.AnalyzerHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
