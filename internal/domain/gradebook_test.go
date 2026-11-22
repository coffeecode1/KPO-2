package domain_test

import (
	"math"
	"testing"

	"KPO-2/internal/domain"
)

func preparedBook(t *testing.T) *domain.Gradebook {
	t.Helper()
	book := domain.NewGradebook()
	if _, err := book.AddStudent("Анна"); err != nil {
		t.Fatal(err)
	}
	if _, err := book.AddCriterion("ДЗ", 10, 2); err != nil {
		t.Fatal(err)
	}
	return book
}

func TestStudentIdentityAndName(t *testing.T) {
	book := domain.NewGradebook()
	if _, err := book.AddStudent(" \t"); err == nil {
		t.Fatal("пустое имя должно отклоняться")
	}
	first, _ := book.AddStudent(" Анна ")
	second, _ := book.AddStudent("Анна")
	if first.ID() == second.ID() || first.Name() != "Анна" {
		t.Fatal("имя обрезается, однофамильцы получают разные ID")
	}
}

func TestCriterionInvariants(t *testing.T) {
	tests := []struct {
		name        string
		max, weight float64
	}{
		{"", 10, 1}, {"ДЗ", 0, 1}, {"ДЗ", -1, 1},
		{"ДЗ", 10, 0}, {"ДЗ", 10, -1},
		{"ДЗ", math.NaN(), 1}, {"ДЗ", 10, math.NaN()},
		{"ДЗ", math.Inf(1), 1}, {"ДЗ", 10, math.Inf(1)},
	}
	for _, test := range tests {
		book := domain.NewGradebook()
		if _, err := book.AddCriterion(test.name, test.max, test.weight); err == nil {
			t.Errorf("принят некорректный критерий: %+v", test)
		}
		if len(book.Rubric()) != 0 {
			t.Fatal("неудачное создание изменило рубрикатор")
		}
	}
}

func TestEvaluationInvariants(t *testing.T) {
	tests := []struct {
		name               string
		student, criterion int
		score              float64
	}{
		{"отрицательный балл", 1, 1, -1},
		{"выше максимума", 1, 1, 12},
		{"несуществующий студент", 9, 1, 5},
		{"несуществующий критерий", 1, 9, 5},
		{"NaN", 1, 1, math.NaN()},
		{"бесконечность", 1, 1, math.Inf(1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			book := preparedBook(t)
			if err := book.SetEvaluation(test.student, test.criterion, test.score); err == nil {
				t.Fatal("некорректная оценка принята")
			}
			if len(book.Evaluations()) != 0 {
				t.Fatal("некорректная оценка сохранилась")
			}
		})
	}
}

func TestEvaluationBoundariesCorrectionAndCancellation(t *testing.T) {
	book := preparedBook(t)
	for _, score := range []float64{0, 10, 7.5} {
		if err := book.SetEvaluation(1, 1, score); err != nil {
			t.Fatal(err)
		}
		if evaluations := book.Evaluations(); len(evaluations) != 1 || evaluations[0].Score().Value() != score {
			t.Fatalf("исправление создало дубль или сохранило не тот балл: %v", evaluations)
		}
	}
	if err := book.SetEvaluation(1, 1, 11); err == nil {
		t.Fatal("ошибочное исправление принято")
	}
	if book.Evaluations()[0].Score().Value() != 7.5 {
		t.Fatal("ошибка стёрла прежнюю оценку")
	}
	if err := book.RemoveEvaluation(1, 1); err != nil {
		t.Fatal(err)
	}
	if len(book.Evaluations()) != 0 {
		t.Fatal("оценка не удалилась")
	}
	if err := book.RemoveEvaluation(1, 1); err == nil {
		t.Fatal("повторная отмена должна сообщить об отсутствии оценки")
	}
}

func TestListsCannotChangeAggregate(t *testing.T) {
	book := preparedBook(t)
	if err := book.SetEvaluation(1, 1, 8); err != nil {
		t.Fatal(err)
	}
	book.Students()[0] = domain.Student{}
	book.Rubric()[0] = domain.Criterion{}
	book.Evaluations()[0] = domain.Evaluation{}
	if book.Students()[0].ID() != 1 || book.Rubric()[0].MaxScore() != 10 || book.Evaluations()[0].Score().Value() != 8 {
		t.Fatal("изменение возвращённого списка повредило агрегат")
	}
}

func TestWeightSumOverflow(t *testing.T) {
	book := domain.NewGradebook()
	if _, err := book.AddCriterion("Первый", 10, math.MaxFloat64); err != nil {
		t.Fatal(err)
	}
	if _, err := book.AddCriterion("Второй", 10, math.MaxFloat64); err == nil {
		t.Fatal("переполнение суммы весов должно отклоняться")
	}
}
