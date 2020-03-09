package main

import (
	"log"
	"net/http"
	"time"

	"gee"
	"middleware"
)

func main() {
	r := gee.New()
	r.Use(middleware.Logger())
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>version 1</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(middlewareV2())
	{
		v2.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>version 2</h1>")
		})

		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}

func middlewareV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
