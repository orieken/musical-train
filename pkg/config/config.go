package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log/slog"
)

// Config holds the application configuration
type Config struct {
	BaseDir string `mapstructure:"base_dir"`
	Timeout int    `mapstructure:"timeout"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		BaseDir: "${HOME}/.master-mold",
		Timeout: 10,
	}
}

// LoadConfig loads the configuration from the specified paths
func LoadConfig(configPaths []string, logger *slog.Logger) (*Config, error) {
	logger.Info("Loading configuration")

	// Set up viper for configuration
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	
	// Add config paths
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	
	v.AutomaticEnv() // Enable environment variable substitution

	// Load the configuration
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a default config
			logger.Info("No config file found, creating default config")

			defaultConfig := DefaultConfig()
			v.Set("base_dir", defaultConfig.BaseDir)
			v.Set("timeout", defaultConfig.Timeout)

			// Ensure the config directory exists
			configDir := filepath.Dir(v.ConfigFileUsed())
			if configDir == "" {
				// If no config file was found, use the first config path
				if len(configPaths) > 0 {
					configDir = configPaths[0]
				} else {
					configDir = "."
				}
			}

			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, errors.Wrap(err, "failed to create config directory")
			}

			v.SetConfigFile(filepath.Join(configDir, "config.toml"))
			if err := v.SafeWriteConfig(); err != nil {
				return nil, errors.Wrap(err, "failed to create default config")
			}
		} else {
			return nil, errors.Wrap(err, "failed to read config file")
		}
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	logger.Info("Configuration loaded", "base_dir", config.BaseDir, "timeout", config.Timeout)
	return &config, nil
}

// GetExpandedBaseDir returns the base directory with environment variables expanded
func GetExpandedBaseDir(config *Config) string {
	return os.ExpandEnv(config.BaseDir)
}

// EnsureBaseDirExists ensures that the base directory exists
func EnsureBaseDirExists(config *Config) error {
	baseDir := GetExpandedBaseDir(config)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return errors.Wrap(err, "failed to create base directory")
		}
	}
	return nil
}