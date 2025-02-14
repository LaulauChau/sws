package client

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LaulauChau/sws/internal/config"
	"github.com/LaulauChau/sws/internal/models"
	"github.com/LaulauChau/sws/pkg/cache"
)

var (
	postTokenURL     = "https://app.sowesign.com/api/portal/authentication/token"
	nextCoursesURL   = "https://app.sowesign.com/api/student-app/future-courses?limit=8"
	currentCourseURL = "https://app.sowesign.com/api/trainer-app/current-courses?limit=1"
)

// SetBaseURLs allows overriding the API endpoints (used for testing)
func SetBaseURLs(tokenURL, coursesURL, currentURL string) {
	postTokenURL = tokenURL
	nextCoursesURL = coursesURL
	currentCourseURL = currentURL
}

type Client struct {
	httpClient *http.Client
	token      string
	config     config.Config
	cache      *cache.Cache[[]models.Course]
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewClient(config config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true, // Often faster for small responses
				ForceAttemptHTTP2:   true,
			},
		},
		config: config,
		cache:  cache.NewCache[[]models.Course](24 * time.Hour),
	}
}

func isTesting() bool {
	return flag.Lookup("test.v") != nil
}

func (c *Client) GetToken() error {
	if c.config.CodeEtablissement == "" || c.config.Identifiant == "" || c.config.PIN == "" {
		return fmt.Errorf("empty credentials provided")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.config.CodeEtablissement + c.config.Identifiant + c.config.PIN))

	req, err := http.NewRequest("POST", postTokenURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "JBAuth "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status code %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if authResp.Token == "" {
		return fmt.Errorf("received empty token from server")
	}

	c.token = "Bearer " + authResp.Token
	if !isTesting() {
		fmt.Println("Successfully obtained token")
	}
	return nil
}

func (c *Client) GetNextCourses() ([]models.Course, error) {
	if courses, ok := c.cache.Get(); ok {
		if !isTesting() {
			fmt.Println("Retrieved courses from cache")
		}
		return courses, nil
	}

	if c.token == "" {
		return nil, fmt.Errorf("no authentication token available")
	}

	req, err := http.NewRequest("GET", nextCoursesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")

	if !isTesting() {
		fmt.Println("Sending request to get next courses...")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned status code %d: %s", resp.StatusCode, string(body))
	}

	var courses []models.Course
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &courses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	c.cache.Set(courses)

	if !isTesting() {
		fmt.Printf("Retrieved %d courses from API and cached\n", len(courses))
	}
	return courses, nil
}
