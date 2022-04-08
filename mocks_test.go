package errors

import "fmt"

const (
	// `0 ~ 100` 为 package "github.com/eachinchung/errors" 保留的错误代码。
	errConfigurationNotValid = iota + 2
	errInvalidJSON
	errEOF
	errLoadConfigFailed
	errNotExt
)

func init() {
	codes[errConfigurationNotValid] = defaultCoder{errConfigurationNotValid, 500, "configuration not valid error"}
	codes[errInvalidJSON] = defaultCoder{errInvalidJSON, 500, "encoding failed due to an error with the data"}
	codes[errEOF] = defaultCoder{errEOF, 500, "end of input"}
	codes[errLoadConfigFailed] = defaultCoder{errLoadConfigFailed, 500, "load configuration file failed"}
	codes[errNotExt] = defaultCoder{errNotExt, 500, ""}
}

func loadConfig() error {
	err := decodeConfig()
	return WithCodef(err, errConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return WithCodef(err, errInvalidJSON, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return WithCodef(err, errEOF, "could not read configuration file")
}
