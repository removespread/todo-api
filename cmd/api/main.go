package main

import (
	"go.uber.org/fx"

	"todo-api/internal/config"
	"todo-api/pkg/logger"
)

func main() {
	fx.New(
		createApp(),
	).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.Provide(
			config.LoadConfig,
			logger.NewLogger,
		),
		fx.Invoke(),
	)
}
