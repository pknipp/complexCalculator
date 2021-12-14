package main

import (
	// "fmt"
	// "io"
	"log"
	"math/cmplx"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

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

func unary (method string, z complex128) (string, complex128) {
	ONE := complex(1., 0.)
	var result complex128
	var message string
	switch method {
		case "Abs":
			result = complex(cmplx.Abs(z), 0.)
		case "Acos":
			result = cmplx.Acos(z)
		case "Acosh":
			result = cmplx.Acosh(z)
		case "Acot":
			if isNonzero(z, &message) {
				result = cmplx.Atan(ONE/z)
			}
		case "Acoth":
			if isNonzero(z, &message) {
				result = cmplx.Atanh(ONE/z)
			}
		case "Acsc":
			if isNonzero(z, &message) {
				result = cmplx.Asin(ONE/z)
			}
		case "Acsch":
			if isNonzero(z, &message) {
				result = cmplx.Asinh(ONE/z)
			}
		case "Asec":
			if isNonzero(z, &message) {
				result = cmplx.Acos(ONE/z)
			}
		case "Asech":
			if isNonzero(z, &message) {
				result = cmplx.Acosh(ONE/z)
			}
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
			if isNonzero(z, &message) {
				result = ONE/cmplx.Tan(z)
			}
		case "Coth":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Tanh(z)
			}
		case "Csc":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Sin(z)
			}
		case "Csch":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Sinh(z)
			}
		case "Exp":
			result = cmplx.Exp(z)
		case "Imag":
			result = complex(imag(z), 0.)
		case "Log":
			if isNonzero(z, &message) {
				result = cmplx.Log(z)
			}
		case "Log10":
			if isNonzero(z, &message) {
				result = cmplx.Log10(z)
			}
		case "Log2":
			if isNonzero(z, &message) {
				result = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
			}
		case "Phase":
			result = complex(cmplx.Phase(z), 0.)
		case "Real":
			result = complex(real(z), 0.)
		case "Sec":
			result = ONE/cmplx.Cos(z)
		case "Sech":
			result = ONE/cmplx.Cosh(z)
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
		default:
			message = "There exists no such function by this name.  Check spelling and capitalization."
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

func parseExpression (expression string) (string, complex128) {
	ZERO, TEN := complex(0., 0.), complex(10., 0.)
	message := ""
	// Following pre-processing line is needed if/when this code is tested in a non-server configuration.
	expression = doRegExp(expression)
	getNumber := func(expression string) (string, complex128, string){
		var val complex128
		message := ""
		if len(expression) == 0 {
			return "Your expression truncates prematurely.", val, expression
		}
		leadingChar := expression[0:1]
		if leadingChar == "(" {
			var nExpression int
			// remove leading parenthesis
			expression = expression[1:]
			message, nExpression = findSize(expression)
			if len(message) != 0 {
				return message, ZERO, ""
			}
			// recursive call to evalulate what is in parentheses
			message, val = parseExpression(expression[0:nExpression])
			if len(message) != 0 {
				return message, ZERO, ""
			}
			// From expression remove trailing parenthesis and stuff preceding it.
			expression = expression[nExpression + 1:]
			return message, val, expression
		} else if leadingChar == "i" {
			return message, complex(0, 1), expression[1:]
			// A letter triggers that we are looking at start of a unary function name.
		} else if isLetter(leadingChar[0]) {
			// If leadingChar is lower-case, convert it to uppercase to facilitate comparison w/our list of unaries.
			if (leadingChar[0] > 96) {
				leadingChar = string(leadingChar[0] - 32)
			}
			expression = expression[1:]
			if len(expression) == 0 {
				return "This unary function invocation ends prematurely.", ZERO, ""
			}
			// If the 2nd character's a letter, this is an invocation of a unary function.
			if isLetter(expression[0]) {
				method := leadingChar
				// We seek an open paren, which signifies start of argument (& end of method name)
				for expression[0:1] != "(" {
					method += expression[0: 1]
					expression = expression[1:]
					if len(expression) == 0 {
						return "The argument of this unary function seems nonexistent.", ZERO, ""
					}
				}
				var nExpression int
				// Remove leading parenthesis
				expression = expression[1:]
				message, nExpression = findSize(expression)
				var arg complex128
				if len(message) != 0 {
					return message, ZERO, ""
				}
				message, arg = parseExpression(expression[0: nExpression])
				if len(message) != 0 {
					return message, ZERO, ""
				}
				message, val = unary(method, arg)
				return message, val, expression[nExpression + 1:]
				// If not a unary, the user is representing scientific notation
			} else if leadingChar[0] == 'E' {
				message = "Your scientific notation (the start of " + leadingChar + expression + ") is improperly formatted."
				p := 1
				for len(expression) >= p {
					if z := expression[0:p]; z != "+" && z != "-" {
						if num, err := strconv.ParseInt(z, 10, 64); err != nil {
							break
						} else {
							val = cmplx.Pow(TEN, complex(float64(num), 0.))
							message = ""
						}
					}
					p++
				}
				return message, val, expression[p - 1:]
			}
		} else {
			// The following'll change only if strconv.ParseFloat ever returns no error, below.
			message = "The string '" + expression + "' does not evaluate to a number."
			p := 1
			for len(expression) >= p {
				// If implied multiplication is detected ...
				if  z := expression[0:p]; expression[p - 1: p] == "(" {
					// ... insert a "*" symbol.
					expression = expression[0:p - 1] + "*" + expression[p - 1:]
					break
				} else if !(z == "." || z == "-" || z == "-.") {
					if num, err := strconv.ParseFloat(z, 64); err != nil {
						break
					} else {
						val = complex(num, 0.)
						message = ""
					}
				}
				p++
			}
			return message, val, expression[p - 1:]
		}
		return "Could not parse " + leadingChar + expression, ZERO, ""
	}
	type opNum struct {
		op string
		num complex128
	}
	if len(expression) > 0 {
		if expression[0:1] == "+" {
			expression = expression[1:]
		}
	}
	var z, num complex128
	message, z, expression = getNumber(expression)
	if len(message) != 0 {
		return message, ZERO
	}
	PRECEDENCE := map[string]int{"+": 0, "-": 0, "*": 1, "/": 1, "^": 2}
	OPS := "+-*/^"
	pairs := []opNum{}
	for len(expression) > 0 {
		op := expression[0:1]
		if strings.Contains(OPS, op) {
			expression = expression[1:]
		} else {
			op = "*"
		}
		if message, num, expression = getNumber(expression); len(message) != 0 {
			return message, ZERO
		} else {
			pairs = append(pairs, opNum{op, num})
		}
	}
	for len(pairs) > 0 {
		index := 0
		for len(pairs) > index {
			if index < len(pairs) - 1 && PRECEDENCE[pairs[index].op] < PRECEDENCE[pairs[index + 1].op] {
				index++
			} else {
				var z1, result complex128
				if index == 0 {
					z1 = z
				} else {
					z1 = pairs[index - 1].num
				}
				message, result = binary(z1, pairs[index].op, pairs[index].num)
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
	return message, z
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
	// Use the following when testing the app in a non-server configuration.
	// expression := "1+2id(3-4id(5+6i))"
	// fmt.Println(handler(expression))
}
