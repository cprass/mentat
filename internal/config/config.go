package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flagType string

const (
	flagTypeString flagType = "string"
	flagTypeBool   flagType = "bool"
)

type flagDef struct {
	flagType    flagType
	name        string
	configKey   string
	shorthand   string
	usage       string
	defaultFunc func() (any, error)
}

var flags = []flagDef{
	{
		flagType:  flagTypeString,
		name:      "config",
		shorthand: "c",
		usage:     "path to config file",
		defaultFunc: func() (any, error) {
			var configDir string
			if runtime.GOOS == "windows" {
				dir, err := os.UserConfigDir()
				if err != nil {
					return nil, err
				}
				configDir = dir
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return nil, err
				}
				configDir = filepath.Join(homeDir, ".config")
			}
			return filepath.Join(configDir, "mentat", "config.yaml"), nil
		},
	},
	{
		flagType:  flagTypeString,
		name:      "vault",
		shorthand: "v",
		usage:     "path to vault dir",
		defaultFunc: func() (any, error) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			return filepath.Join(homeDir, ".mentat"), nil
		},
	},
	{
		flagType:  flagTypeBool,
		name:      "syncEnabled",
		usage:     "automatic vault syncing enabled",
		configKey: "sync.enabled",
	},
	{
		flagType:  flagTypeString,
		name:      "gitRepo",
		usage:     "vault sync: git repository URL",
		configKey: "sync.git.repo",
	},
	{
		flagType: flagTypeString,
		name:     "gitBranch",
		usage:    "vault sync: git branch",
		defaultFunc: func() (any, error) {
			return "main", nil
		},
		configKey: "sync.git.branch",
	},
}

func Vault() (string, error) {
	vaultDir := viper.GetString("vault")

	// Expand $HOME in vault directory
	if strings.HasPrefix(vaultDir, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		vaultDir = path.Join(homeDir, vaultDir[5:])
	}

	return vaultDir, nil
}

func SyncEnabled() bool {
	return viper.GetBool("sync.enabled")
}

type GitConfig struct {
	Url    string
	Branch string
}

func Git() GitConfig {
	return GitConfig{
		Url:    viper.GetString("sync.git.repo"),
		Branch: viper.GetString("sync.git.branch"),
	}
}

// Call this before anything else in the root.go init function
func InitFlags(rootCmd *cobra.Command) {
	for _, config := range flags {
		switch config.flagType {
		case flagTypeBool:
			// load the default value as bool
			defaultValue := false
			if config.defaultFunc != nil {
				result, err := config.defaultFunc()
				if err == nil {
					if boolValue, ok := result.(bool); ok {
						defaultValue = boolValue
					}
				}
			}
			// add the default value to usage information
			rootCmd.PersistentFlags().BoolP(config.name, config.shorthand, defaultValue, config.usage)
		case flagTypeString:
			// load the default value as string
			defaultValue := ""
			if config.defaultFunc != nil {
				result, err := config.defaultFunc()
				if err == nil {
					if strValue, ok := result.(string); ok {
						defaultValue = strValue
					}
				}
			}
			// add the default value to usage information
			rootCmd.PersistentFlags().StringP(config.name, config.shorthand, defaultValue, config.usage)
		}
		configKey := config.configKey
		if configKey == "" {
			configKey = config.name
		}
		viper.BindPFlag(configKey, rootCmd.PersistentFlags().Lookup(config.name))
	}

	// load config file after flags are set up
	cobra.OnInitialize(func() {
		configFile := viper.GetString("config")
		if configFile != "" {
			viper.SetConfigFile(configFile)
		}

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return
			}
			if os.IsNotExist(err) {
				return
			}
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	})
}
