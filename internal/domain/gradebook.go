package domain

import "fmt"

// Gradebook – ведомость курса, все изменения проходят через нее
type Gradebook struct {
	students    []Student
	rubric      []Criterion
	evaluations []Evaluation
}

func NewGradebook() *Gradebook { return &Gradebook{} }

// Возвращаем копии списков
func (b *Gradebook) Students() []Student {
	return append([]Student(nil), b.students...)
}

func (b *Gradebook) Rubric() []Criterion {
	return append([]Criterion(nil), b.rubric...)
}

func (b *Gradebook) Evaluations() []Evaluation {
	return append([]Evaluation(nil), b.evaluations...)
}

func (b *Gradebook) AddStudent(name string) (Student, error) {
	name, err := cleanName(name)
	if err != nil {
		return Student{}, err
	}
	student := Student{id: len(b.students) + 1, name: name}
	b.students = append(b.students, student)
	return student, nil
}

func (b *Gradebook) AddCriterion(name string, maxScore, weight float64) (Criterion, error) {
	name, err := cleanName(name)
	if err != nil {
		return Criterion{}, err
	}
	if !finite(maxScore) || maxScore <= 0 {
		return Criterion{}, fmt.Errorf("максимум должен быть обычным положительным числом")
	}
	if !finite(weight) || weight <= 0 {
		return Criterion{}, fmt.Errorf("вес должен быть обычным положительным числом")
	}
	// проверка на переполнение
	totalWeight := weight
	for _, c := range b.rubric {
		totalWeight += c.Weight()
	}
	if !finite(totalWeight) {
		return Criterion{}, fmt.Errorf("сумма весов слишком велика")
	}
	criterion := Criterion{id: len(b.rubric) + 1, name: name, maxScore: maxScore, weight: weight}
	b.rubric = append(b.rubric, criterion)
	return criterion, nil
}

func (b *Gradebook) Student(id int) (Student, error) {
	for _, student := range b.students {
		if student.ID() == id {
			return student, nil
		}
	}
	return Student{}, fmt.Errorf("студент %d не найден", id)
}

// SetEvaluation выставляет или исправляет оценку, при ошибке старая сохраняется
func (b *Gradebook) SetEvaluation(studentID, criterionID int, value float64) error {
	if _, err := b.Student(studentID); err != nil {
		return err
	}
	evaluation, err := newEvaluation(studentID, criterionID, value, b.rubric)
	if err != nil {
		return err
	}
	for i, old := range b.evaluations {
		if old.StudentID() == studentID && old.CriterionID() == criterionID {
			b.evaluations[i] = evaluation
			return nil
		}
	}
	b.evaluations = append(b.evaluations, evaluation)
	return nil
}

func (b *Gradebook) RemoveEvaluation(studentID, criterionID int) error {
	for i, evaluation := range b.evaluations {
		if evaluation.StudentID() == studentID && evaluation.CriterionID() == criterionID {
			b.evaluations = append(b.evaluations[:i], b.evaluations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("оценка для студента %d по критерию %d не найдена", studentID, criterionID)
}
