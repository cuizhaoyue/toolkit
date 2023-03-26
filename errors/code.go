package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder = defaultCoder{
		BCode:    1,
		HttpCode: http.StatusInternalServerError,
		Ext:      "An internal server error occurred",
		Ref:      "",
	}
)

type Coder interface {
	// HTTP status that should be used for the associated error code
	HTTPStatus() int
	// External (user) faced error text.
	String() string
	// Reference returns the detail documents for user.
	Reference() string
	// Code returns the code of the coder.
	Code() int
}

type defaultCoder struct {
	// Bcode refers to the integer business code of the ErrCode.
	BCode int
	// HttpCode refers to the http status that should be used for the associated error code.
	HttpCode int
	// External (user) faced error text.
	Ext string
	// Ref specify the reference document.
	Ref string
}

// Code returns the integer business code fo the coder.
func (coder defaultCoder) Code() int {
	return coder.BCode
}

// String implements stringer. String returns the external error message, if any.
func (coder defaultCoder) String() string {
	return coder.Ext
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder defaultCoder) HTTPStatus() int {
	if coder.HttpCode == 0 {
		return 500
	}

	return coder.HttpCode
}

// Reference returns the reference document.
func (coder defaultCoder) Reference() string {
	return coder.Ref
}

// codes contains a map of error codes to metadata.
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register register a user define error code.
// It will overrid the existing code.
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved as `unknownCode` error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister register a user define error code.
// It will panic when the same Code already exist.
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by as `unknownCode` error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}
	codes[coder.Code()] = coder
}

// ParseCoder parse any error into *withCode.
// nil error will return nil directly.
// None withStack error will be parsed as ErrUnknown.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode reports whether any error in err's chain contains the given error code.
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}
