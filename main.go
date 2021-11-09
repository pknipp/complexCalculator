package main

import (
	"log"
	"net/http"
	"os"
	"math/cmplex"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func calculate(z1 complex128, op string, z2 complex128) complex128 {
	var result complex128
	switch op {
	case "+":
		result = z1 + z2
	case "-":
		result = z1 - z2
	case "*":
		result = z1 * z2
	case "/":
		result = z1 / z2
	case "^":
		result = cmplx.Pow(z1, z2)
	}
	return result
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}
