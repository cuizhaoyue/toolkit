package errors

import "fmt"

const (
	// Error codes below 1000 are reserved.
	ConfigurationNotValid = iota + 1000
	ErrInvalidJson
	ErrEOF
	ErrLoadConfigFailed
)

func init() {
	Register(defaultCoder{ConfigurationNotValid, 500, "ConfigurationNotValid error", ""})
	Register(defaultCoder{ErrInvalidJson, 500, "Data is not valid JSON", ""})
	Register(defaultCoder{ErrEOF, 500, "End of input", ""})
	Register(defaultCoder{ErrLoadConfigFailed, 500, "Load configuration file failed", ""})
}

func loadConfig() error {
	err := decodeConfig()
	return WrapC(err, ConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return WrapC(err, ErrInvalidJson, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return WrapC(err, ErrEOF, "could not read configuration file")
}
