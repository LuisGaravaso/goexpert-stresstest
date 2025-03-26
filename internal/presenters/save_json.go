package presenters

import (
	"encoding/json"
	"os"
	"stresstest/internal/usecase/run"
)

func SaveReportAsJSON(report run.RunOutputDTO, filePath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if ".json" != filePath[len(filePath)-5:] {
		filePath += ".json"
	}
	return os.WriteFile(filePath, data, 0644)
}
