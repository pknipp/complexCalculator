package main

import (
	"strings"
	"strconv"
)

func ParseExpression (expression string) (string, complex128) {
	ZERO := complex(0., 0.)
	message := ""
	// Following pre-processing line is needed if/when this code is tested in a non-server configuration.
	expression = doRegExp(expression)
	getNumber := func(expression string) (string, complex128, string){
		// ZERO := complex(0., 0.)
		message := ""
		var val complex128
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
			// A capital letter triggers that we are looking at start of a unary function name.
		} else if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", leadingChar) {
			method := leadingChar
			expression = expression[1:]
			if len(expression) == 0 {
				return "This unary function invocation ends prematurely.", 0, expression
			}
			// We seek an open paren, which signifies start of argument (& end of method name)
			for expression[0:1] != "(" {
				method += expression[0: 1]
				expression = expression[1:]
				if len(expression) == 0 {
					return "The argument of this unary function seems nonexistent.", 0, expression
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
	var z complex128
	message, z, expression = getNumber(expression)
	if len(message) != 0 {
		return message, ZERO
	}
	PRECEDENCE := map[string]int{"+": 0, "-": 0, "*": 1, "/": 1, "^": 2}
	OPS := "+-*/^"
	pairs := []opNum{}
	var num complex128
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
				var z1 complex128
				if index == 0 {
					z1 = z
				} else {
					z1 = pairs[index - 1].num
				}
				var result complex128
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
