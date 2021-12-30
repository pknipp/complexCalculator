package main

import (
	"math/cmplx"
	"math"
	"regexp"
	"strconv"
	// _ "github.com/heroku/x/hmetrics/onload"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

// The following enables easy toggling of package between CLI version (for testing) and web version.
var isWebVersion bool = true

func isLetter(char byte) bool {
	if char >= 'A' && char <= 'Z' {
		return true
	} else if char >= 'a' && char <= 'z' {
		return true
	}
	return false
}

func isNonzero(z complex128, m *string) bool {
	isZero := z == complex(0., 0.);
	if isZero {
		*m = "A singularity exists in this expression."
	}
	return !isZero
}

func binary(z1 complex128, op string, z2 complex128) (string, complex128) {
	var result complex128
	var message string
	switch op {
		case "+":
			result = z1 + z2
		case "-":
			result = z1 - z2
		case "*":
			result = z1 * z2
		case "/":
			if isNonzero(z2, &message) {
				result = z1 / z2
			}
		case "^":
			if real(z2) > 0 || isNonzero(z1, &message) {
				result = cmplx.Pow(z1, z2)
			}
		default:
			message = "The operation " + op + " is unknown." // I think that this'll never be hit, because of my use of OPS.
	}
	return message, result
}

func findSize (expression string) (string, int) {
	nParen := 1
	for nExpression := 0; nExpression < len(expression); nExpression++ {
		if char := expression[nExpression: nExpression + 1]; char == "(" {
			nParen++
		} else if char == ")" {
			nParen--
		}
		// Closing parenthesis has been found.
		if nParen == 0 {
			return "", nExpression
		}
	}
	return "No closing parenthesis was found for the following string: '" + expression + "'.", 0
}

// I don't think that this function'll ever fail.
func doRegExp(expression string) string {
	expression = regexp.MustCompile(" ").ReplaceAllString(expression, "")
	expression = regexp.MustCompile("j").ReplaceAllString(expression, "i")
	expression = regexp.MustCompile(`\*\*`).ReplaceAllString(expression, "^")
	expression = regexp.MustCompile("div").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile("DIV").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile(`[dD]`).ReplaceAllString(expression, "/")
	return expression
}

func handler(expression string) string {
	// expression = expression[1:] This was used when I used r.URL.path
	message, result := parseExpression(expression)
	if len(message) != 0 {
		return "ERROR: " + message
	}
	realPart := strconv.FormatFloat(real(result), 'f', -1, 64)
	imagPart := strconv.FormatFloat(math.Abs(imag(result)), 'f', -1, 64)
	var resultString string
	if real(result) != 0 {
		resultString = realPart
	}
	if real(result) != 0 && imag(result) != 0 {
		sign := " + "
		if imag(result) < 0 {
			sign = " - "
		}
		resultString += sign
	}
	if imag(result) != 0 {
		if real(result) == 0 && imag(result) < 0 {
			resultString += " - "
		}
		if math.Abs(imag(result)) != 1. {
			resultString += imagPart
		}
		resultString += "i"
	}
	if real(result) == 0 && imag(result) == 0 {
		resultString = "0"
	}
	return resultString
}

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
