package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

// The following enables easy toggling of package between CLI version (for testing) and web version.
var isWebVersion bool = true

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
			c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
				"baseUrl": "https://complex-calculator-2nd-f1d5dad4b63b.herokuapp.com",
				"frag1": "/(Acos(1.2i(3+i))+4d(5-6i))**(7+8i)",
				"frag2": "/3mds+4mds**2*5s",
			})
		})
		expressionText := "your expression"
		resultText := "numerical value"
		router.GET("/:expression", func(c *gin.Context) {
			expression := doRegExp(c.Param("expression"))
			resultValue, posUnits, negUnits := handler(expression)
			c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"expressionText": expressionText,
					"expressionValue": expression,
					"resultText": resultText,
					"resultValue": resultValue,
					"posUnits": posUnits,
					"negUnits": negUnits,
			})
		})
		router.GET("/api/:expression", func(c *gin.Context) {
			expression := doRegExp(c.Param("expression"))
			resultValue, _, _ := handler(expression)
			resultString := "{\"" + expressionText + "\": " + expression + ", \"" + resultText + "\": " + resultValue + "}"
			c.String(http.StatusOK, resultString)
		})
		router.NoRoute(func(c *gin.Context) {
		    c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Remember that you must type 'd' or 'D' instead of '/' in order to divide one number by another."})
		})
		router.Run(":" + port)
	} else {
		expression := "-(4-i)Sin(1.57(1))"
		result, message := parseExpression(expression)
		fmt.Println(result, message)
		// fmt.Println(handler(expression))
	}
}
