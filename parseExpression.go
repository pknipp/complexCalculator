package main

import (
	"math/cmplx"
	"strconv"
	"strings"
)

type quantityType struct{
	val complex128
	units map[string]complex128
}

func sliceContains (slice []string, char string) bool {
	for _, thisChar := range slice {
		if thisChar == char {
			return true
		}
	}
	return false
}

func parseExpression (expression string) (quantityType, string) {
	TEN := complex(10., 0.)
	message := ""
	// Following pre-processing line is needed if/when this code is tested in a non-server configuration.
	expression = doRegExp(expression)
	getNumber := func(expression *string) (quantityType, string){
		var val quantityType
		message := ""
		if len(*expression) == 0 {
			return val, "Your expression truncates prematurely."
		}
		leadingChar := (*expression)[0:1]
		if leadingChar == "(" {
			// remove leading parenthesis
			*expression = (*expression)[1:]
			var nExpression int
			nExpression, message = findSize(*expression)
			if len(message) != 0 {
				return val, message
			}
			// recursive call to evalulate what is in parentheses
			val, message = parseExpression((*expression)[0:nExpression])
			if len(message) != 0 {
				return val, message
			}
			// From expression remove trailing parenthesis and stuff preceding it.
			*expression = (*expression)[nExpression + 1:]
			return val, message
		} else if leadingChar == "i" {
			*expression = (*expression)[1:]
			return quantityType{val: complex(0, 1), units: nil}, message
		} else if len(*expression) > 2 && (*expression)[0:2] == "mol" {
				*expression = (*expression)[3:]
				units := map[string]complex128{"mol": complex(1., 0.)}
				return quantityType{val: complex(1, 0), units: units}, message
		} else if sliceContains(unitSlice, leadingChar) {
			*expression = (*expression)[1:]
			units := map[string]complex128{}
			units[leadingChar] = complex(1., 0.)
			return quantityType{val: complex(1, 0), units: units}, message
		} else if len(*expression) > 1 && (*expression)[0:2] == "kg" {
			*expression = (*expression)[2:]
			units := map[string]complex128{"kg": complex(1., 0.)}
			return quantityType{val: complex(1, 0), units: units}, message
		} else if isLetter(leadingChar[0]) {
			// A letter here triggers that we are looking at either start of a unary function name, or E-notation
			// If leadingChar is lower-case, convert it to uppercase to facilitate comparison w/our list of unaries.
			leadingChar = strings.ToUpper(leadingChar)
			*expression = (*expression)[1:]
			if len(*expression) == 0 {
				return val, "This unary function invocation ends prematurely."
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
						return val, "This unary function (" + method + ") does not seem to have an argument."
					}
				}
				var nExpression int
				// Remove leading parenthesis
				*expression = (*expression)[1:]
				nExpression, message = findSize(*expression)
				var arg quantityType
				if len(message) != 0 {
					return val, message
				}
				// recursive call, for argument of unary
				arg, message = parseExpression((*expression)[0: nExpression])
				if len(message) != 0 {
					return val, message
				}
				quantity, message := unary(method, arg)
				// Trim argument of unary from beginning of expression
				*expression = (*expression)[nExpression + 1:]
				return quantityType{val: quantity.val, units: quantity.units}, message
			} else if leadingChar[0] == 'E' {
				// If expression is not a unary, the user is representing scientific notation with an "E"
				message = "Your scientific notation (the start of " + leadingChar + *expression + ") is improperly formatted."
				p := 1
				for len(*expression) >= p {
					if z := (*expression)[0:p]; z != "+" && z != "-" {
						if num, err := strconv.ParseInt(z, 10, 64); err != nil {
							break
						} else {
							val = quantityType{val: cmplx.Pow(TEN, complex(float64(num), 0.)), units: nil}
							message = ""
						}
					}
					p++
				}
				*expression = (*expression)[p - 1:]
				return val, message
			}
		} else {
			// The following'll change only if strconv.ParseFloat ever returns no error, below.
			message = "The string '" + *expression + "' does not evaluate to a number."
			p := 1
			for len(*expression) >= p {
				// If implied multiplication is detected ...
				if  z := (*expression)[0:p]; (*expression)[p - 1: p] == "(" {
					// ... insert a "*" symbol.
					*expression = (*expression)[0:p - 1] + "*" + (*expression)[p - 1:]
					break
				} else if !(z == "." || z == "-" || z == "-.") {
					if num, err := strconv.ParseFloat(z, 64); err != nil {
						break
					} else {
						val = quantityType{val: complex(num, 0.), units: nil}
						message = ""
					}
				}
				p++
			}
			*expression = (*expression)[p - 1:]
			return val, message
		}
		return val, "Could not parse " + leadingChar + *expression
	}
	// struct fields consist of binary operation and 2nd number of the pair
	type opNum struct {
		op string
		num quantityType
	}

	if len(expression) > 0 {
		// leading "+" may be trimmed thoughtlessly
		if expression[0:1] == "+" {
			expression = expression[1:]
		}
	}
	// val is lead quantity, and nextVal is any of the following ones
	val, nextVal := quantityType{val: 0, units: nil}, quantityType{val: 0, units: nil}
	pairs := []opNum{}
	// trim&store leading number from expression
	val, message = getNumber(&expression)
	if len(message) != 0 {
		return val, message
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
		if nextVal, message = getNumber(&expression); len(message) != 0 {
			return nextVal, message
		} else {
			pairs = append(pairs, opNum{op, nextVal})
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
				var v1, result quantityType
				if index == 0 {
					v1 = val
				} else {
					v1 = pairs[index - 1].num
				}
				result, message = binary(v1, pairs[index].op, pairs[index].num)
				// mutate the values of z and pairs (reducing the length of the latter by one)
				if index == 0 {
					val = result
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
	return val, message
}
