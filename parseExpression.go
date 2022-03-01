package main

import (
	"math/cmplx"
	"strconv"
	"strings"
)

type unitType struct{
	name string
	power complex128
}

type quantityType struct{
	val complex128
	units []unitType
}

func newONE () complex128 {
	return complex(1., 0.)
}

func newUnits (index int) []unitType {
	var units = []unitType{}
	for _, unit := range UNITS {
		units = append(units, unitType{unit, complex(0., 0.)})
	}
	if index > 0 {
		units[index].power = complex(1., 0.)
	}
	return units
}

// func newQuantity (val complex128, index int) quantityType {
	// return quantityType{val, newUnits(index)}
// }

func sliceContains (slice []string, char string) int {
	for k, thisChar := range slice {
		if thisChar == char {
			return k
		}
	}
	return -1
}

func parseExpression (expression string) (quantityType, string) {
	TEN := complex(10., 0.)
	ZERO := complex(0., 0.)
	ONE := newONE()
	message := ""
	// Following pre-processing line is needed if/when this code is tested in a non-server configuration.
	expression = doRegExp(expression)
	getQuantity := func(expression *string) (quantityType, string){
		message := ""
		quantity := quantityType{ZERO, newUnits(-1)}
		if len(*expression) == 0 {
			return quantity, "Your expression truncates prematurely."
		}
		leadingChar := (*expression)[0:1]
		if leadingChar == "(" {
			// remove leading parenthesis
			*expression = (*expression)[1:]
			var nExpression int
			nExpression, message = findSize(*expression)
			if len(message) != 0 {
				return quantity, message
			}
			// recursive call to evalulate what is in parentheses
			quantity, message = parseExpression((*expression)[0:nExpression])
			if len(message) != 0 {
				return quantity, message
			}
			// From expression remove trailing parenthesis and stuff preceding it.
			*expression = (*expression)[nExpression + 1:]
			return quantity, message
		} else if (leadingChar == "i" || leadingChar == "j") {
			*expression = (*expression)[1:]
			quantity.val = complex(0, 1)
			return quantity, message
		} else if len(*expression) > 2 && (*expression)[0:3] == "mol" {
			k := sliceContains(UNITS, "mol")
			*expression = (*expression)[3:]
			return quantityType{ONE, newUnits(k)}, message
		} else if k := sliceContains(UNITS, leadingChar); k >= 0 {
			*expression = (*expression)[1:]
			units := newUnits(k)
			return quantityType{ONE, units}, message
		} else if len(*expression) > 1 && (*expression)[0:2] == "kg" {
			k:= sliceContains(UNITS, "kg")
			*expression = (*expression)[2:]
			units := newUnits(k)
			return quantityType{ONE, units}, message
		} else if isLetter(leadingChar[0]) {
			// A letter here triggers that we are looking at either start of a unary function name, or E-notation
			// If leadingChar is lower-case, convert it to uppercase to facilitate comparison w/our list of unaries.
			leadingChar = strings.ToUpper(leadingChar)
			*expression = (*expression)[1:]
			if len(*expression) == 0 {
				return quantity, "This unary function invocation ends prematurely."
			}
			if isLetter((*expression)[0]) {
				// If the 2nd character's a letter, this is an invocation of a unary function.
				method := leadingChar
				// We seek an open paren, which signifies start of argument (& end of method name)
				for (*expression)[0:1] != "(" {
					// Append letter to name of method, and trim that letter from beginning of expression.
					method += strings.ToLower((*expression)[0: 1])
					*expression = (*expression)[1:]
					if len(*expression) == 0 {
						return quantity, "This unary function (" + method + ") does not seem to have an argument."
					}
				}
				var nExpression int
				// Remove leading parenthesis
				*expression = (*expression)[1:]
				nExpression, message = findSize(*expression)
				var arg quantityType
				if len(message) != 0 {
					return quantity, message
				}
				// recursive call, for argument of unary
				arg, message = parseExpression((*expression)[0: nExpression])
				if len(message) != 0 {
					return quantity, message
				}
				// Trim argument of unary from beginning of expression
				*expression = (*expression)[nExpression + 1:]
				return unary(method, arg)
			} else if leadingChar[0] == 'E' {
				// If expression is not a unary, the user is representing scientific notation with an "E"
				message = "Your scientific notation (the start of " + leadingChar + *expression + ") is improperly formatted."
				p := 1
				for len(*expression) >= p {
					if z := (*expression)[0:p]; z != "+" && z != "-" {
						if num, err := strconv.ParseInt(z, 10, 64); err != nil {
							break
						} else {
							quantity.val = cmplx.Pow(TEN, complex(float64(num), 0.))
							message = ""
						}
					}
					p++
				}
				*expression = (*expression)[p - 1:]
				return quantity, message
			}
		} else {
			foundValue := false
			// The following'll change only if strconv.ParseFloat ever returns no error, below.
			message = "The string '" + *expression + "' does not evaluate to a number."
			p := 1
			z := ""
			for len(*expression) >= p {
				z = (*expression)[0:p]
				if !(z == "." || z == "-" || z == "-.") {
					if num, err := strconv.ParseFloat(z, 64); err != nil {
						break
					} else {
						quantity.val = complex(num, 0.)
						foundValue = true
						message = ""
					}
				}
				p++
			}
			// example of edge cases trapped by following: -sin(x) or -(x+1)**2
			if (z[:1] == "-" && p == 2 && len(*expression) > 1) {
				quantity.val = complex(-1., 0.)
				foundValue = true
			}
			if foundValue {
				*expression = (*expression)[p - 1:]
				return quantity, ""
			}
		}
		return quantity, "Could not parse " + leadingChar + *expression
	}
	// struct fields consist of binary operation and 2nd number of the pair
	type opQuant struct {
		op string
		quantity quantityType
	}

	if len(expression) > 0 {
		// leading "+" may be trimmed thoughtlessly
		if expression[0:1] == "+" {
			expression = expression[1:]
		}
	}
	pairs := []opQuant{}
	// trim&store leading number from expression
	quantity, message := getQuantity(&expression)
	if len(message) != 0 {
		return quantity, message
	}
	PRECEDENCE := map[string]int{"+": 0, "-": 0, "*": 1, "/": 1, "^": 2}
	OPS := "+-*/^"
	// loop thru the expression, while trimming off (and storing in "pairs" slice) operation/number pairs
	for len(expression) > 0 {
		op := expression[0:1]
		if strings.Contains(OPS, op) {
			expression = expression[1:]
		} else {
			// It must be implied multiplication, so overwrite value of op.
			op = "*"
		}
		if nextQuantity, message := getQuantity(&expression); len(message) != 0 {
			return nextQuantity, message
		} else {
			pairs = append(pairs, opQuant{op, nextQuantity})
		}
	}
	// loop thru "pairs" slice, evaluating operations in order of their precedence
	for len(pairs) > 0 {
		index := 0
		for len(pairs) > index {
			if index < len(pairs) - 1 && PRECEDENCE[pairs[index].op] < PRECEDENCE[pairs[index + 1].op] {
				// postpone this operation because of its lower prececence
				index++
			} else {
				// perform this operation NOW
				var q1, result quantityType
				if index == 0 {
					q1 = quantity
				} else {
					q1 = pairs[index - 1].quantity
				}
				result, message = binary(q1, pairs[index].op, pairs[index].quantity)
				// mutate the values of z and pairs (reducing the length of the latter by one)
				if index == 0 {
					quantity = result
					pairs = pairs[1:]
				} else {
					pairs[index - 1].quantity = result
					pairs = append(pairs[0: index], pairs[index + 1:]...)
				}
				// Start another loop thru the expression, ISO high-precedence operations.
				index = 0
			}
		}
	}
	return quantity, message
}
