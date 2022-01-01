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

func binary(z1 quantityType, op string, z2 quantityType) (quantityType, string) {
	var result quantityType
	var units unitType
	var message string
	var ok bool
	haveSameUnits := func(z1, z2 quantityType) (bool, string) {
		for unit, power := range z1.units {
			if power != z1.units[unit] {
				return false, "You are adding/subtracting quantities w/different units."
			}
		}
		return true, ""
	}
	switch op {
		case "+":
			ok, message = haveSameUnits(z1, z2)
			if ok {
				result = quantityType{val: z1.val + z2.val, units: z1.units}
			}
		case "-":
			ok, message = haveSameUnits(z1, z2)
			if ok {
				result = quantityType{val: z1.val - z2.val, units: z1.units}
			}
		case "*":
			for unit, power := range z1.units {
				units[unit] = power + z2.units[unit]
			}
			result = quantityType{val: z1.val * z2.val, units: units}
		case "/":
			for unit, power := range z1.units {
				units[unit] = power - z2.units[unit]
			}
			if isNonzero(z2.val, &message) {
				result = quantityType{val: z1.val / z2.val, units: units}
			}
		case "^":
			for _, power := range z2.units {
				if power != 0 {
					return result, "An exponent cannot have units."
				}
			}
			if real(z2.val) > 0 || isNonzero(z1.val, &message) {
				result = quantityType{val: cmplx.Pow(z1.val, z2.val), units: units}
			}
		default:
			// I think that this'll never be hit, because of my use of OPS in an outer function.
			message = "The operation " + op + " is unknown."
	}
	return result, message
}

func findSize (expression string) (int, string) {
	nParen := 1 // leading (open)paren has been found, in calling function
	for nExpression := 0; nExpression < len(expression); nExpression++ {
		if char := expression[nExpression: nExpression + 1]; char == "(" {
			nParen++
		} else if char == ")" {
			nParen--
		}
		if nParen == 0 {
			// Closing parenthesis has been found.
			return nExpression, ""
		}
	}
	return 0, "No closing parenthesis was found for this string: (" + expression
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
	result, message := parseExpression(expression)
	if len(message) != 0 {
		return "ERROR: " + message
	}
	realPart := strconv.FormatFloat(real(result.val), 'f', -1, 64)
	imagPart := strconv.FormatFloat(math.Abs(imag(result.val)), 'f', -1, 64)
	var resultString string
	if real(result.val) != 0 {
		resultString = realPart
	}
	if real(result.val) != 0 && imag(result.val) != 0 {
		sign := " + "
		if imag(result.val) < 0 {
			sign = " - "
		}
		resultString += sign
	}
	if imag(result.val) != 0 {
		if real(result.val) == 0 && imag(result.val) < 0 {
			resultString += " - "
		}
		if math.Abs(imag(result.val)) != 1. {
			resultString += imagPart
		}
		resultString += "i"
	}
	if real(result.val) == 0 && imag(result.val) == 0 {
		resultString = "0"
	}
	return resultString
}
