package appManageAPI

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/platform9/appctl/pkg/constants"
)

var dummyConfig = Config{
	IDToken:   "",
	ExpiresAt: time.Now(),
}

var (
	dummyConfigFile string = "config.json"

	dummyBaseDir        string = os.TempDir()
	dummyConfigDir      string = filepath.Join(dummyBaseDir, ".config", "pf9")
	dummyConfigFilePath string = filepath.Join(dummyConfigDir, dummyConfigFile)
)

func logError(t *testing.T, err error) {
	t.Errorf("failed with error: %s\n", err.Error())
}

func TestCreateDirectoryIfNotExist(t *testing.T) {
	if err := createDirectoryIfNotExist(dummyConfigDir); err != nil {
		logError(t, err)
	}
}

func TestCreateConfig(t *testing.T) {
	if err := createConfig(dummyConfig, dummyConfigFilePath); err != nil {
		logError(t, err)
	}
}

func TestLoadConfig(t *testing.T) {
	var loadedConfig *Config
	var err error
	if loadedConfig, err = loadConfig(dummyConfigFilePath); err != nil {
		logError(t, err)
	}
	if dummyConfig.ExpiresAt.Format(constants.UTCClusterTimeStamp) == loadedConfig.ExpiresAt.Format(constants.UTCClusterTimeStamp) ||
		dummyConfig.IDToken != loadedConfig.IDToken {
		errorMessage := fmt.Errorf(`configs donot match
			loadedConfig:
			%+v
			dummyConfig:
			%+v

		`, loadedConfig, dummyConfig)
		logError(t, errorMessage)
	}
}

func TestRemoveConfig(t *testing.T) {
	if err := removeConfig(dummyConfigFilePath); err != nil {
		logError(t, err)
	}
}
