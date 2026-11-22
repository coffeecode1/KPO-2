package adapters

import (
	"encoding/csv"
	"os"
	"strconv"

	"KPO-2/internal/application"
)

// CSVExporter сохраняет итоги в CSV-файл
type CSVExporter struct{}

func NewCSVExporter() *CSVExporter { return &CSVExporter{} }

func (e *CSVExporter) Export(path string, results []application.Result) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	records := [][]string{{"student_id", "student_name", "final_score"}}
	for _, result := range results {
		records = append(records, []string{
			strconv.Itoa(result.StudentID),
			result.Name,
			strconv.FormatFloat(result.Final, 'f', 2, 64),
		})
	}
	writeErr := writer.WriteAll(records)
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
