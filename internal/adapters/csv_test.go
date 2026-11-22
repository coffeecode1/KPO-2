package adapters_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"KPO-2/internal/adapters"
	"KPO-2/internal/application"
)

func TestCSVExporterEscapesNamesAndWritesScores(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.csv")
	exporter := adapters.NewCSVExporter()
	results := []application.Result{{StudentID: 1, Name: "Анна, \"Аня\"\nИванова", Final: 8.6}}
	if err := exporter.Export(path, results); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	got, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"student_id", "student_name", "final_score"}, {"1", results[0].Name, "8.60"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CSV = %v; нужен %v", got, want)
	}
	if err := exporter.Export(t.TempDir(), results); err == nil {
		t.Fatal("нельзя записать CSV поверх каталога")
	}
}

func TestMemoryRepositoryKeepsSameAggregate(t *testing.T) {
	repository := adapters.NewMemoryRepository()
	if _, err := repository.Get().AddStudent("Анна"); err != nil {
		t.Fatal(err)
	}
	if len(repository.Get().Students()) != 1 {
		t.Fatal("хранилище потеряло изменение агрегата")
	}
}
