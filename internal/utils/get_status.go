package utils

import "github.com/gochelias/webaudit/internal/models"

func GetStatus(data []models.Issue, icon string) string {
	len := len(data)

	if len > 0 {
		return icon
	}

	return "✅"
}
