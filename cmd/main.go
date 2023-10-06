package main

import (
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/taraslis453/solid-software-test/config"
	"github.com/taraslis453/solid-software-test/pkg/logging"

	"github.com/taraslis453/solid-software-test/internal/app"
)

func main() {
	logger := logging.NewZapLogger("main")

	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		logger.Fatal("failed to read env", "err", err)
	}
	logger.Info("read config", "config", cfg)

	app.Run(&cfg)
}
