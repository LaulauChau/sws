package handler

import (
	"net/http"

	"github.com/LaulauChau/sws/internal/client"
	"github.com/LaulauChau/sws/internal/config"
	"github.com/LaulauChau/sws/web/templates"
)

type WebHandler struct {
	client *client.Client
}

func NewWebHandler(cfg config.Config) *WebHandler {
	return &WebHandler{
		client: client.NewClient(cfg),
	}
}

func (h *WebHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if err := h.client.GetToken(); err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	courses, err := h.client.GetNextCourses()
	if err != nil {
		http.Error(w, "Failed to get courses", http.StatusInternalServerError)
		return
	}

	if err := templates.Index(courses).Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if err := h.client.GetToken(); err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	courses, err := h.client.GetNextCourses()
	if err != nil {
		http.Error(w, "Failed to get courses", http.StatusInternalServerError)
		return
	}

	if err := templates.CoursesTable(courses).Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
