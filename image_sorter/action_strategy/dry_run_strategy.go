package action_strategy

import (
	"ImageManager/i18n"
	"ImageManager/log"
	"encoding/csv"
	"fmt"
	"os"
)

// DryRunStrategy defines a strategy that logs intended sort operations
// to a CSV file without performing any actual file moves.
type DryRunStrategy struct {
	CsvPath   string
	Header    []string
	csvWriter *csv.Writer
	csvFile   *os.File
}

// NewDryRunStrategy creates a new DryRunStrategy.
func NewDryRunStrategy(csvPath string) *DryRunStrategy {
	return &DryRunStrategy{
		CsvPath: csvPath,
		Header:  []string{i18n.T("SortCsvHeaderSource"), i18n.T("SortCsvHeaderDestination")},
	}
}

// Setup creates the CSV log file and writes the header.
func (s *DryRunStrategy) Setup() error {
	log.LogInfo(i18n.T("SortDryRunActive"))
	var err error
	s.csvFile, err = os.Create(s.CsvPath)
	if err != nil {
		return fmt.Errorf("could not create CSV log file: %w", err)
	}
	s.csvWriter = csv.NewWriter(s.csvFile)
	if s.Header != nil {
		return s.csvWriter.Write(s.Header)
	}
	return nil
}

// Execute logs the planned move to the console and the CSV file.
func (s *DryRunStrategy) Execute(sourcePath, destinationPath string) error {
	fmt.Printf("[DRY RUN] Would move file from %s to %s\n", sourcePath, destinationPath)
	if s.csvWriter != nil {
		row := []string{sourcePath, destinationPath}
		if err := s.csvWriter.Write(row); err != nil {
			log.LogError(fmt.Sprintf("Error writing sort log for %s: %v", sourcePath, err))
		}
	}
	return nil
}

// Teardown flushes the CSV writer and closes the file.
func (s *DryRunStrategy) Teardown() {
	if s.csvWriter != nil {
		s.csvWriter.Flush()
	}
	if s.csvFile != nil {
		log.LogInfo(i18n.T("SortDryRunLogCreated", map[string]interface{}{"Path": s.csvFile.Name()}))
		s.csvFile.Close()
	}
}
