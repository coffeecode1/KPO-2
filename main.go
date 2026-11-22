package main

import (
	"fmt"
	"os"

	"KPO-2/internal/console"
	"KPO-2/internal/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	container, err := di.NewContainer()
	if err != nil {
		return err
	}
	return container.Invoke(func(app *console.App) error {
		return app.Run()
	})
}
