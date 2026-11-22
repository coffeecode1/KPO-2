package application

import (
	"fmt"

	"KPO-2/internal/domain"
)

// Service – фасад, объединяет сценарии приложения для консоли
type Service struct {
	repository GradebookRepository
	calculator *domain.Calculator
	exporter   ResultExporter
}

func NewService(repository GradebookRepository, calculator *domain.Calculator, exporter ResultExporter) *Service {
	return &Service{repository: repository, calculator: calculator, exporter: exporter}
}

func (s *Service) AddStudent(name string) (domain.Student, error) {
	return s.repository.Get().AddStudent(name)
}

func (s *Service) AddCriterion(name string, maxScore, weight float64) (domain.Criterion, error) {
	return s.repository.Get().AddCriterion(name, maxScore, weight)
}

func (s *Service) Students() []domain.Student {
	return s.repository.Get().Students()
}

func (s *Service) Rubric() []domain.Criterion {
	return s.repository.Get().Rubric()
}

func (s *Service) Evaluations() []domain.Evaluation {
	return s.repository.Get().Evaluations()
}

func (s *Service) SetEvaluation(studentID, criterionID int, score float64) error {
	return s.repository.Get().SetEvaluation(studentID, criterionID, score)
}

func (s *Service) RemoveEvaluation(studentID, criterionID int) error {
	return s.repository.Get().RemoveEvaluation(studentID, criterionID)
}

// Results считает итоги всех студентов по выбранному правилу
func (s *Service) Results(strategy domain.FinalStrategy) ([]Result, error) {
	book := s.repository.Get()
	if len(book.Students()) == 0 {
		return nil, fmt.Errorf("сначала добавьте студентов")
	}
	results := make([]Result, 0, len(book.Students()))
	for _, student := range book.Students() {
		final, err := s.calculator.Final(book, student.ID(), strategy)
		if err != nil {
			return nil, err
		}
		results = append(results, Result{StudentID: student.ID(), Name: student.Name(), Final: final})
	}
	return results, nil
}

// Export пересчитывает итоги перед выгрузкой
func (s *Service) Export(path string, strategy domain.FinalStrategy) error {
	results, err := s.Results(strategy)
	if err != nil {
		return err
	}
	return s.exporter.Export(path, results)
}
