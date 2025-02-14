package mock

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/LaulauChau/sws/internal/models"
)

// Server represents a mock HTTP server for testing
type Server struct {
	*httptest.Server
	TokenRequests  int
	CourseRequests int
}

// NewServer creates and returns a new mock server
func NewServer() *Server {
	s := &Server{}

	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers for browser compatibility
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Validate Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		switch r.URL.Path {
		case "/api/portal/authentication/token":
			s.handleTokenRequest(w, r)
		case "/api/student-app/future-courses":
			s.handleCoursesRequest(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))

	return s
}

func (s *Server) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	s.TokenRequests++

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" || auth[:6] != "JBAuth" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if authorization contains empty credentials
	if auth == "JBAuth "+base64.StdEncoding.EncodeToString([]byte("")) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": "mock-token-for-testing-purposes-only",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCoursesRequest(w http.ResponseWriter, r *http.Request) {
	s.CourseRequests++

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth != "Bearer mock-token-for-testing-purposes-only" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	now := time.Now()
	courses := []models.Course{
		{
			ID:    137393,
			Name:  "Innover et entreprendre [XDEV003-CTD / 2425S10-PAR1]",
			Date:  now.Format("2006-01-02"),
			Start: "08:00:00+00:00",
			End:   "12:00:00+00:00",
		},
		{
			ID:    137227,
			Name:  "Innover et entreprendre [XDEV003-CTD / 2425S10-PAR1]",
			Date:  now.AddDate(0, 0, 1).Format("2006-01-02"),
			Start: "13:00:00+00:00",
			End:   "16:30:00+00:00",
		},
	}

	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetBaseURL returns the base URL of the mock server
func (s *Server) GetBaseURL() string {
	return s.URL
}

// Close shuts down the mock server
func (s *Server) Close() {
	s.Server.Close()
}
