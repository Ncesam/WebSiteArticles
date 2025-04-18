package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
)

func GenerateDotEnvExample(cfg *Config) error {
	desc, err := cleanenv.GetDescription(cfg, nil)
	if err != nil {
		return fmt.Errorf("failed to generate description: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to find work dir: %w", err)
	}
	filePath := filepath.Join(cwd, "config", ".example.env")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("can't create .env.example: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(desc)
	if err != nil {
		return fmt.Errorf("can't write .env.example: %w", err)
	}

	fmt.Println(".env.example generated ✅")
	return nil
}
