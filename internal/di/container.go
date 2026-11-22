package di

import (
	"os"

	"go.uber.org/dig"

	"KPO-2/internal/adapters"
	"KPO-2/internal/application"
	"KPO-2/internal/console"
	"KPO-2/internal/domain"
)

// NewContainer регистрирует базу где все лежит, расчёт итогов, экспорт и консоль
func NewContainer() (*dig.Container, error) {
	container := dig.New()
	providers := []any{
		func() application.GradebookRepository { return adapters.NewMemoryRepository() },
		func() application.ResultExporter { return adapters.NewCSVExporter() },
		domain.NewCalculator,
		application.NewService,
		func(service *application.Service) *console.App {
			return console.NewApp(service, os.Stdin, os.Stdout)
		},
	}
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}
	return container, nil
}
