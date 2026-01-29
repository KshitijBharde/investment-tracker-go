package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"

	httphandlers "github.com/KshitijBharde/investment-tracker-go/backend/internal/http"
)

var (
	KITE_CONNECT_API_KEY    string
	KITE_CONNECT_API_SECRET string
	kc                      *kiteconnect.Client
)

func main() {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		panic("Error loading .env file")
	}

	port, exists := os.LookupEnv("BACKEND_PORT")
	if !exists {
		port = "7140"
	}

	KITE_CONNECT_API_KEY, exists = os.LookupEnv("KITE_CONNECT_API_KEY")
	if !exists {
		panic("KITE_CONNECT_API_KEY environment variable not set")
	}
	KITE_CONNECT_API_SECRET, exists = os.LookupEnv("KITE_CONNECT_API_SECRET")
	if !exists {
		panic("KITE_CONNECT_API_SECRET environment variable not set")
	}

	kc = kiteconnect.New(KITE_CONNECT_API_KEY)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Initialize all routes
	httphandlers.InitRoutes(r, kc, KITE_CONNECT_API_SECRET)

	url := fmt.Sprintf("http://localhost:%s", port)
	fmt.Printf("Server started listening on %s\n", url)
	http.ListenAndServe(":"+port, r)
}
