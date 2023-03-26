package core

import (
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the default of the message.
	// This message is suitable to the exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {

	}
}
