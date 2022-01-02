package main

import "math/cmplx"

func areNil (units map[string]complex128, method string, message *string) bool {
	for _, power := range units {
		if power != complex(0., 0.) {
			*message = "The argument of " + method + " must be dimensionless."
			return false
		}
	}
	return true
}

func unary (method string, quantity quantityType) (quantityType, string) {
	z := quantity.val
	argUnits := quantity.units
	units := map[string]complex128{}
	ONE := complex(1., 0.)
	var result complex128
	var message string
	switch method {
		case "Abs":
			result = complex(cmplx.Abs(z), 0.)
			units = argUnits
		case "Acos":
			if areNil(argUnits, method, &message) {
				result = cmplx.Acos(z)
			}
		case "Acosh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Acosh(z)
			}
		case "Acot":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Atan(ONE/z)
				}
			}
		case "Acoth":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Atanh(ONE/z)
				}
			}
		case "Acsc":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Asin(ONE/z)
				}
			}
		case "Acsch":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Asinh(ONE/z)
				}
			}
		case "Asec":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Acos(ONE/z)
				}
			}
		case "Asech":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Acosh(ONE/z)
				}
			}
		case "Asin":
			if areNil(argUnits, method, &message) {
				result = cmplx.Asin(z)
			}
		case "Asinh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Asinh(z)
			}
		case "Atan":
			if areNil(argUnits, method, &message) {
				result = cmplx.Atan(z)
			}
		case "Atanh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Atanh(z)
			}
		case "Conj":
			result = cmplx.Conj(z)
		case "Cos":
			if areNil(argUnits, method, &message) {
				result = cmplx.Cos(z)
			}
		case "Cosh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Cosh(z)
			}
		case "Cot":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = ONE/cmplx.Tan(z)
				}
			}
		case "Coth":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = ONE/cmplx.Tanh(z)
				}
			}
		case "Csc":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = ONE/cmplx.Sin(z)
				}
			}
		case "Csch":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = ONE/cmplx.Sinh(z)
				}
			}
		case "Exp":
			if areNil(argUnits, method, &message) {
				result = cmplx.Exp(z)
			}
		case "Imag":
			result = complex(imag(z), 0.)
		case "Log":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Log(z)
				}
			}
		case "Log10":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Log10(z)
				}
			}
		case "Log2":
			if areNil(argUnits, method, &message) {
				if isNonzero(z, &message) {
					result = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
				}
			}
		case "Phase":
			result = complex(cmplx.Phase(z), 0.)
		case "Real":
			result = complex(real(z), 0.)
		case "Sec":
			if areNil(argUnits, method, &message) {
				result = ONE/cmplx.Cos(z)
			}
		case "Sech":
			if areNil(argUnits, method, &message) {
				result = ONE/cmplx.Cosh(z)
			}
		case "Sin":
			if areNil(argUnits, method, &message) {
				result = cmplx.Sin(z)
			}
		case "Sinh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Sinh(z)
			}
		case "Sqrt":
			result = cmplx.Sqrt(z)
			units = map[string]complex128{}
			for unit, power := range argUnits {
				units[unit] = power / complex(2., 0.)
			}
		case "Tan":
			if areNil(argUnits, method, &message) {
				result = cmplx.Tan(z)
			}
		case "Tanh":
			if areNil(argUnits, method, &message) {
				result = cmplx.Tanh(z)
			}
		default:
			message = method + " is a nonexistent function.  Check your spelling."
	}
	return quantityType{result, units}, message
}
