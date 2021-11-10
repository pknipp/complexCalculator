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

func binary(z1 complex128, op string, z2 complex128) complex128 {
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

func unary(method string, z complex128) complex128 {
	var result complex128
	switch method {
	case "Abs":
		result = complex(cmplx.Abs(z), 0.)
	case "Phase":
		result = complex(cmplx.Phase(z), 0.)
	case "Acos":
		result = cmplx.Acos(z)
	case "Acosh":
		result = cmplx.Acosh(z)
	case "Asin":
		result = cmplx.Asin(z)
	case "Asinh":
		result = cmplx.Asinh(z)
	case "Atan":
		result = cmplx.Atan(z)
	case "Atanh":
		result = cmplx.Atanh(z)
	case "Conj":
		result = cmplx.Conj(z)
	case "Cos":
		result = cmplx.Cos(z)
	case "Cosh":
		result = cmplx.Cosh(z)
	case "Cot":
		result = cmplx.Cot(z)
	case "Coth":
		result = complex(1., 0.)/cmplx.Tanh(z)
	case "Csc":
		result = complex(1., 0.)/cmplx.Sin(z)
	case "Csch":
		result = complex(1., 0.)/cmplx.Sinh(z)
	case "Exp":
		result = cmplx.Exp(z)
	case "Imag":
		result = complex(imag(z), 0.)
	case "Log":
		result = cmplx.Log(z)
	case "Log10":
		result = cmplx.Log10(z)
	case "Log2":
		result = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
	case "Real":
		result = complex(real(z), 0.)
	case "Sec":
		result = complex(1., 0.)/cmplx.Cos(z)
	case "Sech":
		result = complex(1., 0.)/cmplx.Cosh(z)
	case "Sin":
		result = cmplx.Sin(z)
	case "Sinh":
		result = cmplx.Sinh(z)
	case "Sqrt":
		result = cmplx.Sqrt(z)
	case "Tan":
		result = cmplx.Tan(z)
	case "Tanh":
		result = cmplx.Tanh(z)
	}
	// Insert a default case, which should trigger an error message.
	return result
}

func findSize (expression string) int {
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
	return nExpression
}

func doRegExp(expression string) string {
	expression = regexp.MustCompile(" ").ReplaceAllString(expression, "")
	expression = regexp.MustCompile("j").ReplaceAllString(expression, "i")
	expression = regexp.MustCompile(`\*\*`).ReplaceAllString(expression, "^")
	expression = regexp.MustCompile("div").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile("DIV").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile(`[dD]`).ReplaceAllString(expression, "/")
	return expression
}

func parseExpression (expression string) complex128 {
	// Pre-processing is needed here if/when this code is tested in a non-server configuration.
	// expression = doRegExp(expression)
	getNumber := func(expression string) (complex128, string){
		leadingChar := expression[0:1]
		if leadingChar == "(" {
			nExpression := findSize(expression)
			// Recursive call
			return parseExpression(expression[1: nExpression]), expression[nExpression + 1:]
		} else if leadingChar == "i" {
			return complex(0, 1), expression[1:]
		} else if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", leadingChar) {
			method := leadingChar
			expression = expression[1:]
			for expression[0:1] != "(" {
				method += expression[0: 1]
				expression = expression[1:]
			}
			nExpression := findSize(expression)
			arg := parseExpression(expression[1: nExpression])
			return unary(method, arg), expression[nExpression + 1:]
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
				result := binary(z1, pairs[index].op, pairs[index].num)
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
	// expression = expression[1:] This was used when I used r.URL.path
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
		resultString += "i"
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
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// expression := "Sqrt(3+4i)"
	// fmt.Println(parseExpression(expression))
}
