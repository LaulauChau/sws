package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LaulauChau/sws/internal/config"
	"github.com/LaulauChau/sws/internal/models"
)

func TestClient_GetToken(t *testing.T) {
	tests := []struct {
		name       string
		response   AuthResponse
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful token retrieval",
			response: AuthResponse{
				Token: "test-token",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "empty token response",
			response: AuthResponse{
				Token: "",
			},
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "server error",
			response:   AuthResponse{},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request headers
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing content-type header")
				}
				if auth := r.Header.Get("Authorization"); auth == "" {
					t.Error("missing authorization header")
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			// Use test config
			cfg := config.NewTestConfig()
			c := NewClient(cfg)
			c.httpClient = server.Client()

			// Override token URL for testing
			postTokenURL = server.URL

			err := c.GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetNextCourses(t *testing.T) {
	tests := []struct {
		name       string
		courses    []models.Course
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful courses retrieval",
			courses: []models.Course{
				{
					ID:    1,
					Name:  "Test Course",
					Date:  "2025-02-10",
					Start: "08:00:00+00:00",
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "server error",
			courses:    nil,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing content-type header")
				}
				if r.Header.Get("Authorization") == "" {
					t.Error("missing authorization header")
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.courses)
			}))
			defer server.Close()

			// Use test config and initialize client
			cfg := config.NewTestConfig()
			c := NewClient(cfg)
			c.token = "test-token"
			c.httpClient = server.Client()

			// Override courses URL for testing
			nextCoursesURL = server.URL

			got, err := c.GetNextCourses()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextCourses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != len(tt.courses) {
				t.Errorf("GetNextCourses() got = %v courses, want %v", len(got), len(tt.courses))
			}
		})
	}
}

func TestClient_GetNextCourses_NoToken(t *testing.T) {
	cfg := config.NewTestConfig()
	c := NewClient(cfg)

	_, err := c.GetNextCourses()
	if err == nil {
		t.Error("Expected error when calling GetNextCourses without token")
	}
}
