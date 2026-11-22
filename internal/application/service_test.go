package application_test

import (
	"errors"
	"math"
	"testing"

	"KPO-2/internal/adapters"
	"KPO-2/internal/application"
	"KPO-2/internal/domain"
)

type fakeExporter struct {
	path    string
	results []application.Result
	err     error
}

func (e *fakeExporter) Export(path string, results []application.Result) error {
	e.path, e.results = path, results
	return e.err
}

func TestDemoAndExportUseCurrentState(t *testing.T) {
	exporter := &fakeExporter{}
	service := application.NewService(adapters.NewMemoryRepository(), domain.NewCalculator(), exporter)
	if _, err := service.Results(domain.Weighted{}); err == nil {
		t.Fatal("пустая ведомость должна сообщить об отсутствии студентов")
	}
	if err := service.Demo(); err != nil {
		t.Fatal(err)
	}
	if len(service.Students()) != 2 || len(service.Rubric()) != 2 || len(service.Evaluations()) != 4 {
		t.Fatal("демо создало неверный набор данных")
	}
	if err := service.Export("out.csv", domain.Weighted{}); err != nil {
		t.Fatal(err)
	}
	if exporter.path != "out.csv" || len(exporter.results) != 2 || math.Abs(exporter.results[0].Final-8.6) > 1e-9 {
		t.Fatalf("неверный экспорт: %+v", exporter)
	}
	if err := service.SetEvaluation(1, 1, 10); err != nil {
		t.Fatal(err)
	}
	if err := service.RemoveEvaluation(1, 2); err != nil {
		t.Fatal(err)
	}
	if err := service.Export("new.csv", domain.RoundedWeighted{}); err != nil {
		t.Fatal(err)
	}
	if exporter.results[0].Final != 4 {
		t.Fatal("экспорт использует старый итог")
	}
	exporter.err = errors.New("не удалось записать файл")
	if err := service.Export("bad.csv", domain.Weighted{}); !errors.Is(err, exporter.err) {
		t.Fatal("ошибка экспортера потерялась")
	}
}

func TestInvalidResultsDoNotCallExporter(t *testing.T) {
	exporter := &fakeExporter{}
	service := application.NewService(adapters.NewMemoryRepository(), domain.NewCalculator(), exporter)
	if _, err := service.AddStudent("Анна"); err != nil {
		t.Fatal(err)
	}
	if err := service.Export("out.csv", domain.Weighted{}); err == nil {
		t.Fatal("экспорт без рубрикатора не должен выполняться")
	}
	if exporter.path != "" {
		t.Fatal("вызван экспортер после ошибки расчёта")
	}
}
