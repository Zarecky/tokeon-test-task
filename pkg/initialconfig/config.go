package initialconfig

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"tokeon-test-task/internal/config"
	"tokeon-test-task/pkg/log"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
)

type IConfig interface {
	Validate() error
}

// LoadConfig accepts logger to track on config change
func LoadConfig(l log.Logger, mainConfig *config.Config) {
	// Load local config
	if err := loadConfigFromEnv(mainConfig, WithValidation(false)); err != nil {
		l.Fatalf("failed to load local config: %v, ", err)
	}

	// Validate main config
	if err := mainConfig.Validate(); err != nil {
		l.Fatalf("failed to validate local config: %v, ", err)
	}
}

// loadConfigFromEnv - load environment variables from `os env`, `.env` file and pass it to struct.
//
// For local development use `.env` file from root project.
//
// LoadConfigNew also call a `Validate` method.
//
// Example:
//
//	var cfg internalConfig.Config
//	if err := config.LoadConfigNew(&cfg); err != nil {
//		log.Fatalf("could not load configuration: %v, ", err)
//	}
func loadConfigFromEnv(cfg IConfig, opts ...ConfigOption) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return fmt.Errorf("config variable must be a pointer")
	}

	options := ConfigOptions{
		Validation: true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.EnvPath == "" {
		pwdDir, err := os.Getwd()
		if err != nil {
			return err
		}
		options.EnvPath = pwdDir
	}

	aconf := aconfig.Config{
		AllowUnknownFields: true,
		SkipFlags:          true,
		Files:              []string{path.Join(options.EnvPath, ".env")},
		FileDecoders: map[string]aconfig.FileDecoder{
			".env": aconfigdotenv.New(),
		},
	}

	loader := aconfig.LoaderFor(cfg, aconf)
	if err := loader.Load(); err != nil {
		return err
	}

	if !options.Validation {
		return nil
	}

	return cfg.Validate()
}
