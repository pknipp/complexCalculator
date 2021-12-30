package main

import "math/cmplx"

func unary (method string, z complex128) (string, complex128) {
	ONE := complex(1., 0.)
	var result complex128
	var message string
	switch method {
		case "Abs":
			result = complex(cmplx.Abs(z), 0.)
		case "Acos":
			result = cmplx.Acos(z)
		case "Acosh":
			result = cmplx.Acosh(z)
		case "Acot":
			if isNonzero(z, &message) {
				result = cmplx.Atan(ONE/z)
			}
		case "Acoth":
			if isNonzero(z, &message) {
				result = cmplx.Atanh(ONE/z)
			}
		case "Acsc":
			if isNonzero(z, &message) {
				result = cmplx.Asin(ONE/z)
			}
		case "Acsch":
			if isNonzero(z, &message) {
				result = cmplx.Asinh(ONE/z)
			}
		case "Asec":
			if isNonzero(z, &message) {
				result = cmplx.Acos(ONE/z)
			}
		case "Asech":
			if isNonzero(z, &message) {
				result = cmplx.Acosh(ONE/z)
			}
		case "Asin":
			result = cmplx.Asin(z)
		case "Asinh":
			result = cmplx.Asinh(z)
		case "Atan":
			result = cmplx.Atan(z)
		case "Atanh":
			result = cmplx.Atanh(z)
		case "Conj":
			result = cmplx.Conj(z)
		case "Cos":
			result = cmplx.Cos(z)
		case "Cosh":
			result = cmplx.Cosh(z)
		case "Cot":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Tan(z)
			}
		case "Coth":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Tanh(z)
			}
		case "Csc":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Sin(z)
			}
		case "Csch":
			if isNonzero(z, &message) {
				result = ONE/cmplx.Sinh(z)
			}
		case "Exp":
			result = cmplx.Exp(z)
		case "Imag":
			result = complex(imag(z), 0.)
		case "Log":
			if isNonzero(z, &message) {
				result = cmplx.Log(z)
			}
		case "Log10":
			if isNonzero(z, &message) {
				result = cmplx.Log10(z)
			}
		case "Log2":
			if isNonzero(z, &message) {
				result = cmplx.Log(z)/cmplx.Log(complex(2., 0.))
			}
		case "Phase":
			result = complex(cmplx.Phase(z), 0.)
		case "Real":
			result = complex(real(z), 0.)
		case "Sec":
			result = ONE/cmplx.Cos(z)
		case "Sech":
			result = ONE/cmplx.Cosh(z)
		case "Sin":
			result = cmplx.Sin(z)
		case "Sinh":
			result = cmplx.Sinh(z)
		case "Sqrt":
			result = cmplx.Sqrt(z)
		case "Tan":
			result = cmplx.Tan(z)
		case "Tanh":
			result = cmplx.Tanh(z)
		default:
			message = "There exists no such function by this name.  Check spelling and capitalization."
	}
	return message, result
}
