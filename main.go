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

	// "pknipp/parseExpression"
)

func isLetter(char byte) bool {
	if char >= 'A' && char <= 'Z' {
		return true
	} else if char >= 'a' && char <= 'z' {
		return true
	}
	return false
}

func binary(z1 complex128, op string, z2 complex128) (string, complex128) {
	var result complex128
	ZERO := complex(0., 0.)
	message := ""
	pole := "A singularity exists in this expression."
	switch op {
	case "+":
		result = z1 + z2
	case "-":
		result = z1 - z2
	case "*":
		result = z1 * z2
	case "/":
		if z2 == ZERO {
			message = pole
		} else {
			result = z1 / z2
		}
	case "^":
		if z1 == ZERO && real(z2) <= 0 {
			message = pole
		} else {
			result = cmplx.Pow(z1, z2)
		}
	default:
		message = pole
	}
	return message, result
}

func unary(method string, z complex128) (string, complex128) {
	zero, one := complex(0., 0.), complex(1., 0.)
	var result complex128
	message := ""
	pole := "A singularity exists in this expression."
	// scinotation := "You are not implementing scientific notation properly."
	switch method {
	case "Abs":
		result = complex(cmplx.Abs(z), 0.)
	case "Acos":
		result = cmplx.Acos(z)
	case "Acosh":
		result = cmplx.Acosh(z)
	case "Acot":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Atan(one/z)
		}
	case "Acoth":
		if z == zero {
			message = pole
	 	} else {
			result = cmplx.Atanh(one/z)
		}
	case "Acsc":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Asin(one/z)
		}
	case "Acsch":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Asinh(one/z)
		}
	case "Asec":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Acos(one/z)
		}
	case "Asech":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Acosh(one/z)
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
		if z == zero {
			message = pole
		} else {
			result = cmplx.Cot(z)
		}
	case "Coth":
		if z == zero {
			message = pole
		} else {
			result = one/cmplx.Tanh(z)
		}
	case "Csc":
		den := cmplx.Sin(z)
		if den == zero {
			message = pole
		} else {
			result = one/den
		}
	case "Csch":
		den := cmplx.Sinh(z)
		if den == zero {
			message = pole
		} else {
			result = one/den
		}
	case "Exp":
		result = cmplx.Exp(z)
	case "Imag":
		result = complex(imag(z), 0.)
	case "Log":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Log(z)
		}
	case "Log10":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Log10(z)
		}
	case "Log2":
		if z == zero {
			message = pole
		} else {
			result = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
		}
	case "Phase":
		result = complex(cmplx.Phase(z), 0.)
	case "Real":
		result = complex(real(z), 0.)
	case "Sec":
		result = one/cmplx.Cos(z)
	case "Sech":
		result = one/cmplx.Cosh(z)
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
		char := expression[nExpression: nExpression + 1]
		if char == "(" {
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
		// } else if (leadingChar[0] > 96 && leadingChar[0] < 123) || (leadingChar[0] > 64 && leadingChar[0] < 91) {
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
					z := expression[0:p]
					if z != "+" && z != "-" {
						num, err := strconv.ParseInt(z, 10, 64)
						if err != nil {
							break
						}
						val = cmplx.Pow(TEN, complex(float64(num), 0.))
						message = ""
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
					val = complex(num, 0.)
					message = ""
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
		message, num, expression = getNumber(expression)
		if len(message) != 0 {
			return message, ZERO
		}
		pairs = append(pairs, opNum{op, num})
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
	var message, resultString string
	var result complex128
	message, result = parseExpression(expression)
	if len(message) != 0 {
		return "ERROR: " + message
	}
	realPart := strconv.FormatFloat(real(result), 'f', -1, 64)
	imagPart := ""
	// DRY the following with math.abs ASA I figure out how to import it.
	if imag(result) > 0 {
		imagPart = strconv.FormatFloat(imag(result), 'f', -1, 64)
	} else {
		imagPart = strconv.FormatFloat(imag(-result), 'f', -1, 64)
	}
	resultString = ""
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
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// expression := "(1+2id(3-4id(5+6i)))**i"
	// message, resultString := parseExpression(expression)
	// if len(message) == 0 {
		// fmt.Println(resultString)
	// } else {
		// fmt.Println(message)
	// }
}
