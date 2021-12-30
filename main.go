package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

// The following enables easy toggling of package between CLI version (for testing) and web version.
var isWebVersion bool = false

func main() {
	if isWebVersion {
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
		// I opted not to use this version of router, for technical reasons.
		// router := gin.New()
		router := gin.Default()
		router.Use(gin.Logger())
		router.LoadHTMLGlob("templates/*.tmpl.html")
		router.Static("/static", "static")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl.html", nil)
		})
		expressionText := "your expression"
		resultText := "numerical value"
		router.GET("/:expression", func(c *gin.Context) {
			expression := doRegExp(c.Param("expression"))
			c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"expressionText": expressionText,
					"expressionValue": expression,
					"resultText": resultText,
					"resultValue": handler(expression),
			})
		})
		router.GET("/json/:expression", func(c *gin.Context) {
			expression := doRegExp(c.Param("expression"))
			resultString := "{\"" + expressionText + "\": " + expression + ", \"" + resultText + "\": " + handler(expression) + "}"
			c.String(http.StatusOK, resultString)
		})
		router.NoRoute(func(c *gin.Context) {
		    c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Remember that you must type 'd' or 'D' instead of '/' in order to divide one number by another."})
		})
		router.Run(":" + port)
	} else {
		expression := "1+2id(3-4id(5+6i))"
		fmt.Println(handler(expression))
	}
}
