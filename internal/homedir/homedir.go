package homedir

import (
	"errors"
	"os"
)

//nolint:gochecknoglobals //ok
var envKeys = []string{
	"PWD",
	"ALLUSERSAPPDATA",
	"APPDATA",
	"HOME",
}

func Get() (string, error) {
	envHomeDir, exists := tryEnvs(envKeys)

	if exists {
		return envHomeDir, nil
	}

	return "", errors.New("homedir: could not provide homedir. %w")
}

func tryEnvs(envKeys []string) (string, bool) {
	for _, envKey := range envKeys {
		pathToTry, exists := os.LookupEnv(envKey)

		if !exists {
			continue
		}

		isOk := existsDir(pathToTry)

		if isOk {
			return pathToTry, true
		}
	}

	return "", false
}

func existsDir(directory string) bool {
	_, err := os.Stat(directory)

	return err == nil
}
