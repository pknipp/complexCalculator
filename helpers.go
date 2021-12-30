package main

import (
	"math/cmplx"
	"math"
	"regexp"
	"strconv"
)

func isLetter(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func isNonzero(z complex128, m *string) bool {
	isZero := z == complex(0., 0.)
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
	nParen := 1 // leading paren has been found, in calling function
	for nExpression := 0; nExpression < len(expression); nExpression++ {
		if char := expression[nExpression: nExpression + 1]; char == "(" {
			nParen++
		} else if char == ")" {
			nParen--
		}
		if nParen == 0 {
			// Closing parenthesis has been found.
			return "", nExpression
		}
	}
	return "No closing parenthesis was found for the following string: (" + expression, 0
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
