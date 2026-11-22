package adapters

import "KPO-2/internal/domain"

// MemoryRepository хранит ведомость в памяти до завершения программы
type MemoryRepository struct {
	book *domain.Gradebook
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{book: domain.NewGradebook()}
}

func (r *MemoryRepository) Get() *domain.Gradebook { return r.book }
