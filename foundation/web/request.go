package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type validator interface {
	Validate() error
}

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	// gin
	c, ok := r.Context().(*gin.Context)
	if ok {
		return c.Param(key)
	}

	return ""
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
func Decode(c *gin.Context, val any) error {
	if err := c.ShouldBindJSON(val); err != nil {
		// This will catch all the validation issues with the JSON provided.
		return wrapError(err)
	}
	return nil
}
