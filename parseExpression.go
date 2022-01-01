package main

import (
	"math/cmplx"
	"strconv"
	"strings"
)

func parseExpression (expression string) (complex128, string) {
	ZERO, TEN := complex(0., 0.), complex(10., 0.)
	message := ""
	// Following pre-processing line is needed if/when this code is tested in a non-server configuration.
	expression = doRegExp(expression)
	getNumber := func(expression string) (complex128, string, string){
		var val complex128
		message := ""
		if len(expression) == 0 {
			return val, expression, "Your expression truncates prematurely."
		}
		leadingChar := expression[0:1]
		if leadingChar == "(" {
			var nExpression int
			// remove leading parenthesis
			expression = expression[1:]
			nExpression, message = findSize(expression)
			if len(message) != 0 {
				return ZERO, "", message
			}
			// recursive call to evalulate what is in parentheses
			val, message = parseExpression(expression[0:nExpression])
			if len(message) != 0 {
				return ZERO, "", message
			}
			// From expression remove trailing parenthesis and stuff preceding it.
			expression = expression[nExpression + 1:]
			return val, expression, message
		} else if leadingChar == "i" {
			return complex(0, 1), expression[1:], message
			// A letter triggers that we are looking at start of a unary function name.
		} else if isLetter(leadingChar[0]) {
			// If leadingChar is lower-case, convert it to uppercase to facilitate comparison w/our list of unaries.
			leadingChar = strings.ToUpper(leadingChar)
			expression = expression[1:]
			if len(expression) == 0 {
				return ZERO, "", "This unary function invocation ends prematurely."
			}
			// If the 2nd character's a letter, this is an invocation of a unary function.
			if isLetter(expression[0]) {
				method := leadingChar
				// We seek an open paren, which signifies start of argument (& end of method name)
				for expression[0:1] != "(" {
					// Append letter to name of method, and trim that letter from beginning of expression.
					method += strings.ToLower(expression[0: 1])
					expression = expression[1:]
					if len(expression) == 0 {
						return ZERO, "", "The argument of this unary function seems nonexistent."
					}
				}
				var nExpression int
				// Remove leading parenthesis
				expression = expression[1:]
				nExpression, message = findSize(expression)
				var arg complex128
				if len(message) != 0 {
					return ZERO, "", message
				}
				arg, message = parseExpression(expression[0: nExpression])
				if len(message) != 0 {
					return ZERO, "", message
				}
				val, message = unary(method, arg)
				// Trim argument of unary from beginning of expression
				return val, expression[nExpression + 1:], message
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
				return val, expression[p - 1:], message
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
			return val, expression[p - 1:], message
		}
		return ZERO, "", "Could not parse " + leadingChar + expression
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
	z, expression, message = getNumber(expression)
	if len(message) != 0 {
		return ZERO, message
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
		if num, expression, message = getNumber(expression); len(message) != 0 {
			return ZERO, message
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
				result, message = binary(z1, pairs[index].op, pairs[index].num)
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
	return z, message
}
