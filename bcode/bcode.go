package bcode

import (
	"github.com/cuizhaoyue/toolkit/errors"
	"github.com/novalagung/gubrak"
)

// ErrCode implementes `errors`.Coder interface.
type ErrCode struct {
	// Bcode refers to the integer business code of the ErrCode.
	BCode int

	// HttpCode refers to the http status that should be used for the associated error code.
	HttpCode int

	// External (user) faced error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.BCode
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	return coder.HttpCode
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200,400,401,403,404,500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		BCode:    code,
		HttpCode: httpStatus,
		Ext:      message,
		Ref:      reference,
	}

	errors.MustRegister(coder)
}
