package jee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%s]\t%s\t%v", c.Req.Method, c.Req.RequestURI, time.Since(start))
	}
}
