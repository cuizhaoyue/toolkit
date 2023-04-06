package code

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
		panic("http code not in `200, 400, 401, 403, 404, 500`")
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

// func init() {
// 	register(ErrUserNotFound, 404, "User not found")
// 	register(ErrUserAlreadyExist, 400, "User already exist")
// 	register(ErrReachMaxCount, 400, "Secret reach the max count")
// 	register(ErrSecretNotFound, 404, "Secret not found")
// 	register(ErrSuccess, 200, "OK")
// 	register(ErrUnknown, 500, "Internal server error")
// 	register(ErrBind, 400, "Error occurred while binding the request body to the struct")
// 	register(ErrValidation, 400, "Validation failed")
// 	register(ErrTokenInvalid, 401, "Token invalid")
// 	register(ErrDatabase, 500, "Database error")
// 	register(ErrEncrypt, 401, "Error occurred while encrypting the user password")
// 	register(ErrSignatureInvalid, 401, "Signature is invalid")
// 	register(ErrExpired, 401, "Token expired")
// 	register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
// 	register(ErrMissingHeader, 401, "The `Authorization` header was empty")
// 	register(ErrPasswordIncorrect, 401, "Password was incorrect")
// 	register(ErrPermissionDenied, 403, "Permission denied")
// 	register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
// 	register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
// 	register(ErrInvalidJSON, 500, "Data is not valid JSON")
// 	register(ErrEncodingJSON, 500, "JSON data could not be encoded")
// 	register(ErrDecodingJSON, 500, "JSON data could not be decoded")
// 	register(ErrInvalidYaml, 500, "Data is not valid Yaml")
// 	register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
// 	register(ErrDecodingYaml, 500, "Yaml data could not be decoded")
// }
