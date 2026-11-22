package console_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"KPO-2/internal/adapters"
	"KPO-2/internal/application"
	"KPO-2/internal/console"
	"KPO-2/internal/domain"
)

func runConsole(t *testing.T, input string) string {
	t.Helper()
	service := application.NewService(adapters.NewMemoryRepository(), domain.NewCalculator(), adapters.NewCSVExporter())
	var output bytes.Buffer
	if err := console.NewApp(service, strings.NewReader(input), &output).Run(); err != nil {
		t.Fatal(err)
	}
	return output.String()
}

func TestConsoleManualScenarioAndInputErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.csv")
	input := strings.Join([]string{
		"unknown", // Неизвестный пункт не завершает программу.
		"1", "Вера",
		"3", "ДЗ", "10", "1",
		"2", "4",
		"5", "abc", // Некорректный id
		"5", "1", "1", "текст", // Некорректное число.
		"5", "1", "1", "12",
		"5", "1", "1", "7,5", // Запятая в дроби
		"6", "8",
		"5", "1", "1", "9",
		"7", "1", "1",
		"8",
		"5", "1", "1", "8,6",
		"10", "2", "8",
		"9", path,
		"10", "1", "0", "",
	}, "\n")
	output := runConsole(t, input)
	for _, want := range []string{
		"неизвестный пункт", "целое число", "нужно ввести число", "от 0 до 10",
		"Оценка сохранена", "7.50", "Оценка отменена", "0.00", "9.00", "Итоги экспортированы",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("в выводе нет %q", want)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "1,Вера,9.00") {
		t.Fatalf("неверный экспорт из консоли: %s", data)
	}
}

func TestConsoleDemo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "demo.csv")
	output := runConsole(t, "11\n"+path+"\n0\n")
	for _, want := range []string{"Вера Иванова", "Борис Петров", "8.60", "5.40", "7.00", "Итоги экспортированы"} {
		if !strings.Contains(output, want) {
			t.Errorf("в демо нет %q", want)
		}
	}
}

func TestConsoleEOFAndEmptyLists(t *testing.T) {
	for _, input := range []string{"", "1\n", "3\nДЗ\n", "5\n1\n", "9\n", "10\n"} {
		runConsole(t, input)
	}
	output := runConsole(t, "2\n4\n6\n8\n10\n3\n0\n")
	if !strings.Contains(output, "Список пуст") || !strings.Contains(output, "сначала добавьте студентов") || !strings.Contains(output, "выберите 1 или 2") {
		t.Fatal("не обработаны пустое состояние или неверная стратегия")
	}
}
