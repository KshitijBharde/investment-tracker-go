package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var (
	accessToken string // Storing access token in memory for testing
)

func InitRoutes(r *chi.Mux, kc *kiteconnect.Client, apiSecret string) {
	// Health check
	r.Get("/api/hc", handleHealthCheck)

	// Kite Connect Testing Routes
	r.Get("/api/kite/login-url", handleGetLoginURL(kc))
	r.Post("/api/kite/session", handleGenerateSession(kc, apiSecret))
	r.Get("/api/kite/mf/holdings", handleGetMFHoldings(kc))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleGetLoginURL(kc *kiteconnect.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginURL := kc.GetLoginURL()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"login_url": loginURL,
			"note":      "Visit this URL to authorize and get request token. Token is valid for ~2 minutes",
		})
	}
}

func handleGenerateSession(kc *kiteconnect.Client, apiSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestToken := r.URL.Query().Get("request_token")
		if requestToken == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "request_token query parameter required",
			})
			return
		}

		data, err := kc.GenerateSession(requestToken, apiSecret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to generate session: %v", err),
			})
			return
		}

		accessToken = data.AccessToken
		kc.SetAccessToken(accessToken)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      true,
			"access_token": data.AccessToken,
			"user_id":      data.UserID,
			"message":      "Session created successfully. Access token valid for this session",
		})
	}
}

func handleGetMFHoldings(kc *kiteconnect.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if accessToken == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "No active session. Call /api/kite/session first",
			})
			return
		}

		holdings, err := kc.GetMFHoldings()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to get MF holdings: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(holdings)
	}
}
