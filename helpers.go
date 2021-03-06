package main

import (
	"math/cmplx"
	"math"
	"regexp"
	"strconv"
	"fmt"
)

var UNITS = []string{"kg", "m", "s", "K", "mol"}

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

func areSame(units1, units2 []unitType) (bool, string) {
	for k := range UNITS {
		if units1[k].power != units2[k].power {
			return false, "You tried to add/subtract quantities w/different units."
		}
	}
	return true, ""
}

func binary(q1 quantityType, op string, q2 quantityType) (quantityType, string) {
	// value & units fields of "result" will be adjusted differently from those of z1 in each case
	quantity := q1
	if len(quantity.units) == 0 {
		quantity = quantityType{complex(0., 0.), newUnits(-1)}
	}
	var message string
	var ok bool
	switch op {
		case "+":
			ok, message = areSame(q1.units, q2.units)
			if ok {
				quantity.val += q2.val
			}
		case "-":
			ok, message = areSame(q1.units, q2.units)
			if ok {
				quantity.val -= q2.val
			}
		case "*":
			quantity.val *= q2.val
			for k, unit := range q2.units {
				quantity.units[k].power += unit.power
			}
		case "/":
			if isNonzero(q2.val, &message) {
				quantity.val /= q2.val
				for k, unit := range q2.units {
					quantity.units[k].power -= unit.power
				}
			}
		case "^":
			for _, unit := range q2.units {
				if unit.power != 0 {
					return quantity, "An exponent cannot have units."
				}
			}
			for k := range q1.units {
				quantity.units[k].power *= q2.val
			}
			if real(q2.val) > 0 || isNonzero(q1.val, &message) {
				quantity.val = cmplx.Pow(q1.val, q2.val)
			}
		default:
			// I think that this'll never be hit, because of my use of OPS in parseExpression.
			message = "The operation " + op + " is unknown."
	}
	return quantity, message
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
	expression = regexp.MustCompile(`\*\*`).ReplaceAllString(expression, "^")
	expression = regexp.MustCompile("div").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile("DIV").ReplaceAllString(expression, "/")
	expression = regexp.MustCompile(`[dD]`).ReplaceAllString(expression, "/")
	return expression
}

type unitPower struct {
	unit string
	power complex128
}

func handler(expression string) (string, [][2]string, [][2]string) {
	// expression = expression[1:] This was used when I used r.URL.path
	result, message := parseExpression(expression)
	posUnits := [][2]string{}
	negUnits := [][2]string{}
	if len(message) != 0 {
		return "ERROR: " + message, posUnits, negUnits
	}
	realPart := strconv.FormatFloat(real(result.val), 'f', -1, 64)
	imagPart := strconv.FormatFloat(math.Abs(imag(result.val)), 'f', -1, 64)
	for _, unit := range result.units {
		var powString string
		reFloat, imFloat := int(real(unit.power)), int(imag(unit.power))
		if float64(reFloat) == real(unit.power) {
			if math.Abs(real(unit.power)) == 1. {
				powString = ""
			} else {
				powString = strconv.Itoa(int(math.Abs(real(unit.power))))
			}
		} else {
			powString = fmt.Sprintf("%.2f", real(unit.power))
		}
		if imag(unit.power) != 0. {
			if float64(imFloat) == imag(unit.power) {
				powString += "+" + strconv.Itoa(int(math.Abs(imag(unit.power)))) + "i"
			} else {
				powString += "+" + fmt.Sprintf("%.2f", math.Abs(imag(unit.power))) + "i"
			}
		}
		if real(unit.power) > 0 {
			posUnits = append(posUnits, [2]string{unit.name,  powString})
		} else if real(unit.power) < 0 {
			negUnits = append(negUnits, [2]string{unit.name,  powString})
		}
	}
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
	return resultString, posUnits, negUnits
}
