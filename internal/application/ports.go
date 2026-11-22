package application

import "KPO-2/internal/domain"

// GradebookRepository – доступ к ведомости
type GradebookRepository interface {
	Get() *domain.Gradebook
}

// Result – строка отчёта с итогом студента
type Result struct {
	StudentID int
	Name      string
	Final     float64
}

// ResultExporter выгружает итоги во внешнее хранилище
type ResultExporter interface {
	Export(path string, results []Result) error
}
