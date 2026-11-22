package domain

import (
	"fmt"
	"math"
	"strings"
)

// Student – студент, у каждого свой id
type Student struct {
	id   int
	name string
}

func (s Student) ID() int      { return s.id }
func (s Student) Name() string { return s.name }

// Criterion – критерий оценивания с максимальным баллом и весом
type Criterion struct {
	id       int
	name     string
	maxScore float64
	weight   float64
}

func (c Criterion) ID() int           { return c.id }
func (c Criterion) Name() string      { return c.name }
func (c Criterion) MaxScore() float64 { return c.maxScore }
func (c Criterion) Weight() float64   { return c.weight }

// Score – балл без собственного id
type Score struct {
	value float64
}

func NewScore(value float64) (Score, error) {
	if !finite(value) || value < 0 {
		return Score{}, fmt.Errorf("балл должен быть конечным неотрицательным числом")
	}
	return Score{value: value}, nil
}

func (s Score) Value() float64 { return s.value }

// Evaluation – оценка студента по одному критерию
type Evaluation struct {
	studentID   int
	criterionID int
	score       Score
}

func (e Evaluation) StudentID() int   { return e.studentID }
func (e Evaluation) CriterionID() int { return e.criterionID }
func (e Evaluation) Score() Score     { return e.score }

// newEvaluation создаёт оценку, проверяет критерий и возможный балл
func newEvaluation(studentID, criterionID int, value float64, rubric []Criterion) (Evaluation, error) {
	for _, criterion := range rubric {
		if criterion.ID() != criterionID {
			continue
		}
		score, err := NewScore(value)
		if err != nil {
			return Evaluation{}, err
		}
		if score.Value() > criterion.MaxScore() {
			return Evaluation{}, fmt.Errorf("балл должен быть от 0 до %g", criterion.MaxScore())
		}
		return Evaluation{studentID: studentID, criterionID: criterionID, score: score}, nil
	}
	return Evaluation{}, fmt.Errorf("критерий %d отсутствует в рубрикаторе", criterionID)
}

func finite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func cleanName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("название или имя не может быть пустым")
	}
	return name, nil
}
