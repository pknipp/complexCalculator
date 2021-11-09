package main

import (
	// "fmt"
	// "io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

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

func parseExpression (expression string) (complex128) {
	getNumber := func(expression string) (complex128, string){
		leadingChar := expression[0:1]
		if leadingChar == "(" {
			nExpression := 0
			nParen := 1
			for nParen > 0 {
				nExpression++
				nextChar := expression[nExpression: nExpression + 1]
				if nextChar == "(" {
					nParen++
				} else if nextChar == ")" {
					nParen--
				}
			}
			return parseExpression(expression[1: nExpression]), expression[nExpression + 1:]
		} else if leadingChar == "i" {
			return complex(0, 1), expression[1:]
		} else {
			p := 1
			var lastNum complex128
			for len(expression) >= p {
				z := expression[0:p]
				// If implied multiplication is detected ...
				if expression[p - 1: p] == "(" {
					// ... insert a "*" symbol.
					expression = expression[0:p - 1] + "*" + expression[p - 1:]
					break
				} else if !(z == "." || z == "-" || z == "-.") {
					num, err := strconv.ParseFloat(z, 64)
					if err != nil {
						break
					}
					lastNum = complex(num, 0.)
				}
				p++
			}
			return lastNum, expression[p - 1:]
		}
	}
	type opNum struct {
		op string
		num complex128
	}

	if expression[0:1] == "+" {
		expression = expression[1:]
	}
	var z complex128
	z, expression = getNumber(expression)
	precedence := map[string]int{"+": 0, "-": 0, "*": 1, "/": 1, "^": 2}
	ops := "+-*/^"
	pairs := []opNum{}
	var num complex128
	for len(expression) > 0 {
		op := expression[0:1]
		if strings.Contains(ops, op) {
			expression = expression[1:]
		} else {
			op = "*"
		}
		num, expression = getNumber(expression)
		pair := opNum{op, num}
		pairs = append(pairs, pair)
	}
	for len(pairs) > 0 {
		index := 0
		for len(pairs) > index {
			if index < len(pairs) - 1 && precedence[pairs[index].op] < precedence[pairs[index + 1].op] {
				index++
			} else {
				var z1 complex128
				if index == 0 {
					z1 = z
				} else {
					z1 = pairs[index - 1].num
				}
				result := calculate(z1, pairs[index].op, pairs[index].num)
				if index == 0 {
					z = result
					pairs = pairs[1:]
				} else {
					pairs[index - 1].num = result
					pairs = append(pairs[0: index], pairs[index + 1:]...)
				}
				// Start another loop thru the expression, ISO high-precedence operations.
				index = 0
			}
		}
	}
	return z
}

func handler(expression string) string {
	// expression = expression[1:]
	result := parseExpression(expression)
	realPart := strconv.FormatFloat(real(result), 'f', -1, 64)
	imagPart := ""
	// DRY the following with math.abs ASA I figure out how to import it.
	if imag(result) > 0 {
		imagPart = strconv.FormatFloat(imag(result), 'f', -1, 64)
	} else {
		imagPart = strconv.FormatFloat(imag(-result), 'f', -1, 64)
	}
	resultString := ""
	if real(result) != 0 {
		resultString += realPart
	}
	if real(result) != 0 && imag(result) != 0 {
		// DRY the following after finding some sort of "sign" function
		if imag(result) > 0 {
			resultString += " + "
		} else {
			resultString += " - "
		}
	}
	if imag(result) != 0 {
		if real(result) == 0 && imag(result) < 0 {
			resultString += " - "
		}
		// DRY the following after figuring out how to import math.abs
		if imag(result) != 1 && imag(result) != -1 {
			resultString += imagPart
		}
		resultString += " i"
	}
	if real(result) == 0 && imag(result) == 0 {
		resultString = "0"
	}
	return resultString
}

// func handlerOld(w http.ResponseWriter, r*http.Request) {
	// io.WriteString(w, "numerical value of the expression above = ")
	// expression := r.URL.Path
	// if expression != "/favicon.ico" {
		// resultString := handler(expression)
		// io.WriteString(w, resultString)
	// }
// }

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
//
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/:expression", func(c *gin.Context) {
		expression := c.Param("expression")
		expression = regexp.MustCompile(" ").ReplaceAllString(expression, "")
		expression = regexp.MustCompile("j").ReplaceAllString(expression, "i")
		expression = regexp.MustCompile(`\*\*`).ReplaceAllString(expression, "^")
		expression = regexp.MustCompile("div").ReplaceAllString(expression, "/")
		expression = regexp.MustCompile("DIV").ReplaceAllString(expression, "/")
		expression = regexp.MustCompile(`[dD]`).ReplaceAllString(expression, "/")
		c.String(http.StatusOK, "your expression = " + expression + "\n")
		resultString := handler(expression)
		c.String(http.StatusOK, "numerical value = " + resultString)
	})
	// http.ListenAndServe(":8000", nil)
	router.Run(":" + port)
}
