package main

import "math/cmplx"

func areNone (units []unitType, method string, message *string) bool {
	for _, unit := range units {
		if unit.power != complex(0., 0.) {
			*message = "The argument of " + method + " must be dimensionless."
			return false
		}
	}
	return true
}

func unary (method string, quantity quantityType) (quantityType, string) {
	ONE := newONE()
	// A few cases (Abs, Re, Imag, Cong, Sqrt) will overwrite these units.
	units := newUnits(-1)
	z := quantity.val
	var val complex128
	var message string
	switch method {
		case "Abs":
			val = complex(cmplx.Abs(z), 0.)
			units = quantity.units
		case "Acos":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Acos(z)
			}
		case "Acosh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Acosh(z)
			}
		case "Acot":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Atan(ONE/z)
				}
			}
		case "Acoth":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Atanh(ONE/z)
				}
			}
		case "Acsc":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Asin(ONE/z)
				}
			}
		case "Acsch":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Asinh(ONE/z)
				}
			}
		case "Asec":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Acos(ONE/z)
				}
			}
		case "Asech":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Acosh(ONE/z)
				}
			}
		case "Asin":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Asin(z)
			}
		case "Asinh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Asinh(z)
			}
		case "Atan":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Atan(z)
			}
		case "Atanh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Atanh(z)
			}
		case "Conj":
			val = cmplx.Conj(z)
			units = quantity.units
		case "Cos":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Cos(z)
			}
		case "Cosh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Cosh(z)
			}
		case "Cot":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = ONE/cmplx.Tan(z)
				}
			}
		case "Coth":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = ONE/cmplx.Tanh(z)
				}
			}
		case "Csc":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = ONE/cmplx.Sin(z)
				}
			}
		case "Csch":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = ONE/cmplx.Sinh(z)
				}
			}
		case "Exp":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Exp(z)
			}
		case "Imag":
			val = complex(imag(z), 0.)
			units = quantity.units
		case "Log":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Log(z)
				}
			}
		case "Log10":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Log10(z)
				}
			}
		case "Log2":
			if areNone(quantity.units, method, &message) {
				if isNonzero(z, &message) {
					val = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
				}
			}
		case "Phase":
			val = complex(cmplx.Phase(z), 0.)
		case "Real":
			val = complex(real(z), 0.)
			units = quantity.units
		case "Sec":
			if areNone(quantity.units, method, &message) {
				val = ONE/cmplx.Cos(z)
			}
		case "Sech":
			if areNone(quantity.units, method, &message) {
				val = ONE/cmplx.Cosh(z)
			}
		case "Sin":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Sin(z)
			}
		case "Sinh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Sinh(z)
			}
		case "Sqrt":
			val = cmplx.Sqrt(z)
			units = quantity.units
			for k := range quantity.units {
				units[k].power /= complex(2., 0.)
			}
		case "Tan":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Tan(z)
			}
		case "Tanh":
			if areNone(quantity.units, method, &message) {
				val = cmplx.Tanh(z)
			}
		default:
			message = method + " is a nonexistent function.  Check your spelling."
	}
	return quantityType{val, units}, message
}
