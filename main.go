package main

import (
	"github.com/joiller/jee/jee"
	"net/http"
)

func main() {
	e := jee.New()
	e.Use(jee.Logger(), jee.Recovery())
	e.GET("/nihao/:name", func(c *jee.Context) {
		c.JSON(http.StatusOK, jee.H{
			"name": c.Param("name"),
		})
	})

	e.GET("/panic", func(c *jee.Context) {
		panic("panic error")
	})

	e.Static("/static", "./assets")
	e.Run(":9999")
}
