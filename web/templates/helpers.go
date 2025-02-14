package templates

import (
	"github.com/LaulauChau/sws/internal/models"
	"github.com/LaulauChau/sws/internal/service"
)

func generateCode(course models.Course) string {
	_, _, _, code := service.GenerateFixedCode(course)
	return code
}
