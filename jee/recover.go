package jee

import (
	"errors"
	"fmt"
	"net/http"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Fail(http.StatusInternalServerError, errors.New(fmt.Sprintf("%s", err)))
			}
		}()
		c.Next()
	}
}
