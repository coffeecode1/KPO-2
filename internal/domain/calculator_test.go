package domain_test

import (
	"math"
	"testing"

	"KPO-2/internal/domain"
)

func checkFinal(t *testing.T, book *domain.Gradebook, strategy domain.FinalStrategy, want float64) {
	t.Helper()
	got, err := domain.NewCalculator().Final(book, 1, strategy)
	if err != nil || math.Abs(got-want) > 1e-9 {
		t.Fatalf("итог = %v, ошибка = %v; нужен %v", got, err, want)
	}
}

func TestFinalUsesWeightsMaximaAndCurrentEvaluations(t *testing.T) {
	book := preparedBook(t) // ДЗ: максимум 10, вес 2.
	if _, err := book.AddCriterion("Экзамен", 20, 3); err != nil {
		t.Fatal(err)
	}
	if _, err := book.AddStudent("Другой студент"); err != nil {
		t.Fatal(err)
	}
	for _, grade := range []struct {
		student, criterion int
		score              float64
	}{{1, 1, 8}, {1, 2, 18}, {2, 1, 10}} {
		if err := book.SetEvaluation(grade.student, grade.criterion, grade.score); err != nil {
			t.Fatal(err)
		}
	}
	checkFinal(t, book, domain.Weighted{}, 8.6)
	checkFinal(t, book, domain.RoundedWeighted{}, 9)
	if err := book.SetEvaluation(1, 2, 10); err != nil {
		t.Fatal(err)
	}
	checkFinal(t, book, domain.Weighted{}, 6.2)
	if err := book.RemoveEvaluation(1, 2); err != nil {
		t.Fatal(err)
	}
	checkFinal(t, book, domain.Weighted{}, 3.2)
	if _, err := book.AddCriterion("Проект", 10, 5); err != nil {
		t.Fatal(err)
	}
	checkFinal(t, book, domain.Weighted{}, 1.6)
}

func TestFinalEmptyCases(t *testing.T) {
	book := preparedBook(t)
	checkFinal(t, book, domain.Weighted{}, 0)
	calculator := domain.NewCalculator()
	if _, err := calculator.Final(book, 99, domain.Weighted{}); err == nil {
		t.Fatal("неизвестный студент")
	}
	if _, err := calculator.Final(book, 1, nil); err == nil {
		t.Fatal("нет стратегии")
	}
	empty := domain.NewGradebook()
	if _, err := empty.AddStudent("Анна"); err != nil {
		t.Fatal(err)
	}
	if _, err := calculator.Final(empty, 1, domain.Weighted{}); err == nil {
		t.Fatal("нет критериев")
	}
}

func TestLargeValidNumbersDoNotOverflowFinal(t *testing.T) {
	book := domain.NewGradebook()
	book.AddStudent("Анна")
	book.AddCriterion("ДЗ", math.MaxFloat64, math.MaxFloat64)
	if err := book.SetEvaluation(1, 1, math.MaxFloat64); err != nil {
		t.Fatal(err)
	}
	checkFinal(t, book, domain.Weighted{}, 10)
}
