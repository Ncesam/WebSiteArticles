package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func SetupConfig() (*Config, error) {
	var cfg Config

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't inititalize zap logger: %v", err)
	}
	cfgFile := filepath.Join(cwd, "config", ".env")
	err = godotenv.Load(cfgFile)
	if err != nil {
		return nil, fmt.Errorf(".env not found or failed to load: %v", err)
	}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("environment failed when try scanning: %v", err)
	}
	err = GenerateDotEnvExample(&cfg)
	if err != nil {
		return nil, fmt.Errorf("example environment failed create: %v", err)
	}
	return &cfg, nil
}
