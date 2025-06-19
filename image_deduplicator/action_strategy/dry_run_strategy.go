package action_strategy

import (
	"ImageManager/i18n"
	"ImageManager/log"
	"encoding/csv"
	"fmt"
	"os"
)

// DryRunStrategy defines a strategy that logs potential duplicates to a CSV file
// without performing any file operations.
type DryRunStrategy struct {
	CsvPath   string
	Header    []string // Header for the CSV file.
	csvWriter *csv.Writer
	csvFile   *os.File
}

// Setup creates the CSV log file and writes the header.
func (a *DryRunStrategy) Setup() error {
	log.LogInfo(i18n.T("DryRunActive"))
	var err error
	a.csvFile, err = os.Create(a.CsvPath)
	if err != nil {
		return fmt.Errorf("could not create CSV log file: %w", err)
	}
	a.csvWriter = csv.NewWriter(a.csvFile)
	// Use the provided headers.
	if a.Header != nil {
		return a.csvWriter.Write(a.Header)
	}
	return nil
}

// Execute logs the file to keep and the file to remove to the CSV.
func (a *DryRunStrategy) Execute(toKeep string, toRemove []string) {
	for _, pathToRemove := range toRemove {
		fmt.Printf("[DRY RUN] Potential duplicate found: %s (Original: %s)\n", pathToRemove, toKeep)
		if a.csvWriter != nil {
			row := []string{toKeep, pathToRemove}
			if err := a.csvWriter.Write(row); err != nil {
				log.LogError(fmt.Sprintf("Error writing duplicate log for %s: %v", pathToRemove, err))
			}
		}
	}
}

// Teardown flushes the CSV writer and closes the file.
func (a *DryRunStrategy) Teardown() {
	if a.csvWriter != nil {
		a.csvWriter.Flush()
	}
	if a.csvFile != nil {
		log.LogInfo(i18n.T("DryRunLogCreated", map[string]interface{}{"Path": a.csvFile.Name()}))
		a.csvFile.Close()
	}
}
