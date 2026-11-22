package di_test

import (
	"testing"

	"KPO-2/internal/adapters"
	"KPO-2/internal/application"
	"KPO-2/internal/console"
	"KPO-2/internal/di"
)

func TestContainerBuildsGraphAndAllowsRepositoryReplacement(t *testing.T) {
	container, err := di.NewContainer()
	if err != nil {
		t.Fatal(err)
	}
	repository := adapters.NewMemoryRepository()
	if _, err := repository.Get().AddStudent("Из подменённого репозитория"); err != nil {
		t.Fatal(err)
	}
	if err := container.Decorate(func() application.GradebookRepository { return repository }); err != nil {
		t.Fatal(err)
	}
	err = container.Invoke(func(app *console.App, service *application.Service) {
		if app == nil || len(service.Students()) != 1 || service.Students()[0].Name() != "Из подменённого репозитория" {
			t.Fatal("граф не собрался либо зависимость не подменена")
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}
