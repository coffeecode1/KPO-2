package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"KPO-2/internal/application"
	"KPO-2/internal/domain"
)

// App читает команды пользователя и выводит результаты
type App struct {
	service  *application.Service
	scanner  *bufio.Scanner
	out      io.Writer
	strategy domain.FinalStrategy
}

func NewApp(service *application.Service, in io.Reader, out io.Writer) *App {
	return &App{
		service: service, scanner: bufio.NewScanner(in), out: out,
		strategy: domain.Weighted{},
	}
}

// Run показывает меню до выхода пользователя или окончания ввода
func (a *App) Run() error {
	fmt.Fprintln(a.out, "ПОС – преподавательская оценка студентов")
	fmt.Fprintln(a.out, "Данные хранятся в памяти до выхода, для понимания выберите 11.")
	for {
		fmt.Fprintln(a.out, "\n1. Добавить студента\n2. Список студентов\n3. Добавить критерий\n4. Рубрикатор\n5. Выставить / исправить оценку\n6. Список оценок\n7. Отменить оценку\n8. Итоги и средний балл\n9. Экспорт итогов в CSV\n10. Правило итога\n11. Демо: студенты → критерии → оценки → итог → экспорт\n0. Выход")
		command, err := a.line("Пункт: ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if command == "0" {
			return nil
		}
		err = a.execute(command)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if a.scanner.Err() != nil {
			return a.scanner.Err()
		}
		if err != nil {
			fmt.Fprintln(a.out, "Ошибка:", err)
		}
	}
}

func (a *App) execute(command string) error {
	switch command {
	case "1":
		return a.addStudent()
	case "2":
		a.students()
	case "3":
		return a.addCriterion()
	case "4":
		a.rubric()
	case "5":
		return a.setEvaluation()
	case "6":
		a.evaluations()
	case "7":
		return a.removeEvaluation()
	case "8":
		return a.results()
	case "9":
		return a.export()
	case "10":
		return a.chooseStrategy()
	case "11":
		if err := a.service.Demo(); err != nil {
			return err
		}
		a.students()
		a.rubric()
		a.evaluations()
		if err := a.results(); err != nil {
			return err
		}
		return a.export()
	default:
		return fmt.Errorf("неизвестный пункт меню")
	}
	return nil
}

// line читает строку и при окончании ввода возвращает io.EOF
func (a *App) line(prompt string) (string, error) {
	fmt.Fprint(a.out, prompt)
	if !a.scanner.Scan() {
		if err := a.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(a.scanner.Text()), nil
}

func (a *App) integer(prompt string) (int, error) {
	line, err := a.line(prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("нужно ввести целое число")
	}
	return value, nil
}

// number читает дробное число с точкой или запятой
func (a *App) number(prompt string) (float64, error) {
	line, err := a.line(prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(strings.ReplaceAll(line, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("нужно ввести число")
	}
	return value, nil
}

func (a *App) addStudent() error {
	name, err := a.line("Имя студента: ")
	if err != nil {
		return err
	}
	student, err := a.service.AddStudent(name)
	if err != nil {
		return err
	}
	fmt.Fprintf(a.out, "Добавлен студент #%d: %s\n", student.ID(), student.Name())
	return nil
}

func (a *App) addCriterion() error {
	name, err := a.line("Название критерия: ")
	if err != nil {
		return err
	}
	maxScore, err := a.number("Максимальный балл: ")
	if err != nil {
		return err
	}
	weight, err := a.number("Вес (положительное число, например 0.4): ")
	if err != nil {
		return err
	}
	criterion, err := a.service.AddCriterion(name, maxScore, weight)
	if err != nil {
		return err
	}
	fmt.Fprintf(a.out, "Добавлен критерий #%d: %s\n", criterion.ID(), criterion.Name())
	return nil
}

// evaluationIDs просит id студента и критерия для выставления оценки
func (a *App) evaluationIDs() (int, int, error) {
	studentID, err := a.integer("ID студента: ")
	if err != nil {
		return 0, 0, err
	}
	criterionID, err := a.integer("ID критерия: ")
	return studentID, criterionID, err
}

func (a *App) setEvaluation() error {
	studentID, criterionID, err := a.evaluationIDs()
	if err != nil {
		return err
	}
	score, err := a.number("Балл: ")
	if err != nil {
		return err
	}
	if err := a.service.SetEvaluation(studentID, criterionID, score); err != nil {
		return err
	}
	fmt.Fprintln(a.out, "Оценка сохранена")
	return nil
}

func (a *App) removeEvaluation() error {
	studentID, criterionID, err := a.evaluationIDs()
	if err != nil {
		return err
	}
	if err := a.service.RemoveEvaluation(studentID, criterionID); err != nil {
		return err
	}
	fmt.Fprintln(a.out, "Оценка отменена")
	return nil
}

func (a *App) students() {
	fmt.Fprintln(a.out, "Студенты:")
	for _, student := range a.service.Students() {
		fmt.Fprintf(a.out, "#%d — %s\n", student.ID(), student.Name())
	}
	if len(a.service.Students()) == 0 {
		fmt.Fprintln(a.out, "Список пуст")
	}
}

func (a *App) rubric() {
	fmt.Fprintln(a.out, "Рубрикатор:")
	for _, criterion := range a.service.Rubric() {
		fmt.Fprintf(a.out, "#%d — %s; максимум %g; вес %g\n", criterion.ID(), criterion.Name(), criterion.MaxScore(), criterion.Weight())
	}
	if len(a.service.Rubric()) == 0 {
		fmt.Fprintln(a.out, "Список пуст")
	}
}

func (a *App) evaluations() {
	fmt.Fprintln(a.out, "Оценки:")
	for _, evaluation := range a.service.Evaluations() {
		fmt.Fprintf(a.out, "Студент #%d, критерий #%d: %g\n", evaluation.StudentID(), evaluation.CriterionID(), evaluation.Score().Value())
	}
	if len(a.service.Evaluations()) == 0 {
		fmt.Fprintln(a.out, "Список пуст")
	}
}

// results выводит итоги студентов и средний итог группы
func (a *App) results() error {
	results, err := a.service.Results(a.strategy)
	if err != nil {
		return err
	}
	table := tabwriter.NewWriter(a.out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tСтудент\tИтог / 10")
	total := 0.0
	for _, result := range results {
		fmt.Fprintf(table, "%d\t%s\t%.2f\n", result.StudentID, result.Name, result.Final)
		total += result.Final
	}
	if err := table.Flush(); err != nil {
		return err
	}
	fmt.Fprintf(a.out, "Средний итог: %.2f\n", total/float64(len(results)))
	return nil
}

func (a *App) export() error {
	path, err := a.line("Путь CSV: ")
	if err != nil {
		return err
	}
	if path == "" {
		path = "results.csv"
	}
	if err := a.service.Export(path, a.strategy); err != nil {
		return err
	}
	fmt.Fprintln(a.out, "Итоги экспортированы:", path)
	return nil
}

func (a *App) chooseStrategy() error {
	choice, err := a.line("1 — взвешенный итог; 2 — тот же итог с округлением до целого: ")
	if err != nil {
		return err
	}
	switch choice {
	case "1":
		a.strategy = domain.Weighted{}
	case "2":
		a.strategy = domain.RoundedWeighted{}
	default:
		return fmt.Errorf("выберите 1 или 2")
	}
	fmt.Fprintln(a.out, "Правило изменено для следующих расчётов")
	return nil
}
