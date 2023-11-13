package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/htquangg/microservices-poc/pkg/constants"

	"github.com/spf13/viper"
)

func LoadConfig(cfg interface{}) (interface{}, error) {
	var configPath string

	env := os.Getenv(constants.AppEnv)
	if env == "" {
		env = constants.Dev
	}
	fmt.Println(os.Getenv(constants.ConfigPath))

	configPathFromEnv := viper.Get(constants.ConfigPath)
	if configPathFromEnv != nil {
		configPath = configPathFromEnv.(string)
	} else {
		// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		d, err := getConfigRootPath()
		if err != nil {
			return nil, err
		}
		configPath = d
	}

	// https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Yaml)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getConfigRootPath() (string, error) {
	// Get the current working directory
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fmt.Printf("Current working directory is: %s\n", wd)

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")

	return configPath, nil
}
