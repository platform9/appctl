package appManageAPI

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var dummyConfig = Config{
	IDToken:   "wizK8eib75MNuw==",
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

func logSubTestFail(t *testing.T, subTest string) {

}

func SubTestCreateDirectoryIfNotExist(t *testing.T) {
	if err := createDirectoryIfNotExist(dummyConfigDir); err != nil {
		logError(t, err)
	}
}

func SubTestCreateConfig(t *testing.T) {
	if err := createConfig(dummyConfig, dummyConfigFilePath); err != nil {
		logError(t, err)
	}
}

func SubTestLoadConfig(t *testing.T) {
	var loadedConfig *Config
	var err error
	if loadedConfig, err = loadConfig(dummyConfigFilePath); err != nil {
		logError(t, err)
	}
	if dummyConfig.ExpiresAt.UnixNano() != loadedConfig.ExpiresAt.UnixNano() ||
		dummyConfig.IDToken != loadedConfig.IDToken {
		errorMessage := fmt.Errorf(`configs do not match
			loadedConfig:
			%+v
			dummyConfig:
			%+v

		`, loadedConfig, dummyConfig)
		logError(t, errorMessage)
	}
}

func SubTestRemoveConfig(t *testing.T) {
	if err := removeConfig(dummyConfigFilePath); err != nil {
		logError(t, err)
	}
}

func TestCreateLoadRemoveConfig(t *testing.T) {
	if !t.Run("SubTestCreateDirectoryIfNotExist", SubTestCreateDirectoryIfNotExist) {
		logSubTestFail(t, t.Name())
	}
	if !t.Run("SubTestCreateConfig", SubTestCreateConfig) {
		logSubTestFail(t, t.Name())
	}
	if !t.Run("SubTestLoadConfig", SubTestLoadConfig) {
		logSubTestFail(t, t.Name())
	}
	if !t.Run("SubTestRemoveConfig", SubTestRemoveConfig) {
		logSubTestFail(t, t.Name())
	}
	t.Cleanup(func() {
		removeConfig(dummyConfigFilePath)
	})
}
