package report

import (
	"encoding/json"
	"os"

	"github.com/gochelias/webaudit/internal/models"
)

func New() models.Report {
	return models.Report{}
}

func Save(path string, r *models.Report) error {
	data, err := json.MarshalIndent(r, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
