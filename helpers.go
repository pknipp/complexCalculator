package main

import (
	"math/cmplx"
	"math"
	"regexp"
	"strconv"
	"fmt"
)

var unitSlice = []string{"kg", "m", "s"}

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
	var message string
	var ok bool
	haveSameUnits := func(z1, z2 quantityType) (bool, string) {
		for _, unit := range unitSlice {
			power1, in1 := z1.units[unit]
			power2, in2 := z2.units[unit]
			if in1 == in2 {
				if power1 == power2 {
					return true, ""
				}
			}
			return false, "You are adding/subtracting quantities w/different units."
		}
		return true, ""
	}
	units := map[string]complex128{}
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
			for _, unit := range unitSlice {
				if power, found := z1.units[unit]; found {
					units[unit] = power
				}
				if power, found := z2.units[unit]; found {
					units[unit] += power
				}
			}
			result = quantityType{val: z1.val * z2.val, units: units}
		case "/":
			for _, unit := range unitSlice {
				if power, found := z1.units[unit]; found {
					units[unit] = power
				}
				if power, found := z2.units[unit]; found {
					units[unit] -= power
				}
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
			for unit, power := range z1.units {
				units[unit] = power * z2.val
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
	for unit, power := range result.units {
		var reString, imString string
		reFloat, imFloat := int(real(power)), int(imag(power))
		if float64(reFloat) == real(power) {
			reString = strconv.Itoa(reFloat)
		} else {
			reString = fmt.Sprintf("%f", real(power))
		}
		if float64(imFloat) == imag(power) {
			imString = strconv.Itoa(int(math.Abs(imag(power))))
		} else {
			imString = fmt.Sprintf("%f", math.Abs(imag(power)))
		}
		if real(power) > 0 {
			posUnits = append(posUnits, [2]string{unit,  reString + "+" + imString + "i"})
		} else {
			negUnits = append(negUnits, [2]string{unit,  reString + "+" + imString + "i"})
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
	return resultString, posUnits, negUnits //+ unitString
}
