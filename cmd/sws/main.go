package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LaulauChau/sws/internal/config"
	"github.com/LaulauChau/sws/internal/handler"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("\nMake sure you have set up your .env file with the following variables:")
		fmt.Println("SOWESIGN_CODE_ETABLISSEMENT")
		fmt.Println("SOWESIGN_IDENTIFIANT")
		fmt.Println("SOWESIGN_PIN")
		os.Exit(1)
	}

	webHandler := handler.NewWebHandler(cfg)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	http.HandleFunc("/", webHandler.HandleIndex)
	http.HandleFunc("/refresh", webHandler.HandleRefresh)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
