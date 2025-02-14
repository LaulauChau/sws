package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LaulauChau/sws/internal/client"
	"github.com/LaulauChau/sws/internal/config"
	"github.com/LaulauChau/sws/internal/models"
)

type testServer struct {
	*httptest.Server
	// Track requests for verification
	tokenRequests  int
	courseRequests int
}

func newTestServer() *testServer {
	ts := &testServer{}

	ts.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate common headers
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch r.URL.Path {
		case "/api/portal/authentication/token":
			ts.handleTokenRequest(w, r)
		case "/api/student-app/future-courses":
			ts.handleCoursesRequest(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return ts
}

func (ts *testServer) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	ts.tokenRequests++

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Validate JBAuth header
	auth := r.Header.Get("Authorization")
	if auth == "" || auth[:6] != "JBAuth" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": "test-integration-token",
	})
}

func (ts *testServer) handleCoursesRequest(w http.ResponseWriter, r *http.Request) {
	ts.courseRequests++

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Validate Bearer token
	auth := r.Header.Get("Authorization")
	if auth != "Bearer test-integration-token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Return realistic test data
	courses := []models.Course{
		{
			ID:    137393,
			Name:  "Test Course 1",
			Date:  time.Now().Format("2006-01-02"),
			Start: "08:00:00+00:00",
			End:   "12:00:00+00:00",
		},
		{
			ID:    137227,
			Name:  "Test Course 2",
			Date:  time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			Start: "13:00:00+00:00",
			End:   "16:30:00+00:00",
		},
	}

	json.NewEncoder(w).Encode(courses)
}

func TestEndToEndFlow(t *testing.T) {
	// Start test server
	ts := newTestServer()
	defer ts.Close()

	// Override API endpoints to use test server
	client.SetBaseURLs(ts.URL+"/api/portal/authentication/token",
		ts.URL+"/api/student-app/future-courses",
		ts.URL+"/api/trainer-app/current-courses")

	// Create client with test config
	cfg := config.Config{
		CodeEtablissement: "test-code",
		Identifiant:       "test-id",
		PIN:               "test-pin",
	}

	c := client.NewClient(cfg)

	// Test authentication flow
	err := c.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if ts.tokenRequests != 1 {
		t.Errorf("Expected 1 token request, got %d", ts.tokenRequests)
	}

	// Test course retrieval
	courses, err := c.GetNextCourses()
	if err != nil {
		t.Fatalf("GetNextCourses failed: %v", err)
	}
	if ts.courseRequests != 1 {
		t.Errorf("Expected 1 course request, got %d", ts.courseRequests)
	}

	// Validate course data
	if len(courses) != 2 {
		t.Errorf("Expected 2 courses, got %d", len(courses))
	}

	// Validate course fields
	for _, course := range courses {
		if course.ID == 0 {
			t.Error("Course ID should not be zero")
		}
		if course.Name == "" {
			t.Error("Course name should not be empty")
		}
		if course.Date == "" {
			t.Error("Course date should not be empty")
		}
		if course.Start == "" {
			t.Error("Course start time should not be empty")
		}
		if course.End == "" {
			t.Error("Course end time should not be empty")
		}
	}
}
