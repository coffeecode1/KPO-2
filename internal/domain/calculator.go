package domain

import (
	"fmt"
	"math"
)

// FinalStrategy – правило получения итога от 0 до 1
type FinalStrategy interface {
	Calculate(normalizedScore float64) float64
}

// Weighted – дробный итог от 0 до 10
type Weighted struct{}

func (Weighted) Calculate(normalizedScore float64) float64 {
	return normalizedScore * 10
}

// RoundedWeighted – итог с округлением до целого, 0.5 округляется вверх
type RoundedWeighted struct{}

func (RoundedWeighted) Calculate(normalizedScore float64) float64 {
	return math.Round(normalizedScore * 10)
}

// Calculator считает итог по текущим оценкам и рубрикатору
type Calculator struct{}

func NewCalculator() *Calculator { return &Calculator{} }

func (c *Calculator) Final(book *Gradebook, studentID int, strategy FinalStrategy) (float64, error) {
	if _, err := book.Student(studentID); err != nil {
		return 0, err
	}
	if len(book.rubric) == 0 {
		return 0, fmt.Errorf("сначала добавьте критерии")
	}
	if strategy == nil {
		return 0, fmt.Errorf("не выбрано правило подсчета")
	}
	totalWeight := 0.0
	for _, criterion := range book.rubric {
		totalWeight += criterion.Weight()
	}
	result := 0.0
	for _, criterion := range book.rubric {
		for _, evaluation := range book.evaluations {
			if evaluation.StudentID() == studentID && evaluation.CriterionID() == criterion.ID() {
				// Чтобы не было переполнения
				result += (evaluation.Score().Value() / criterion.MaxScore()) * (criterion.Weight() / totalWeight)
			}
		}
	}
	// Ограничиваем погрешность float64 чтобы оценка была не больше 10
	return strategy.Calculate(math.Min(result, 1)), nil
}
