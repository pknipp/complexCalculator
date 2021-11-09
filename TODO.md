figure out how to import "math" into handler, which enable DRY-ing of some parts of code.
figure out how to separate functions into separate files
error handling
return json (and html), separately.

unary operations: See https://pkg.go.dev/math/cmplx
real return from complex input (no import required): real, imag ("imag" will conflict w/"i"!)
The following require import of "math/cmplx" & are invoked as, e.g., cmplx.Acos(z).
real return from complex input: Abs, Phase
complex (in&out): Acos,Acosh,Asin,Asinh,Atan,Atanh,Conj,Cos,Cosh,Cot,Exp,Log,Log10,Sin,Sinh,Sqrt,Tan,Tanh
Given one complex input, Polar returns two real numbers: r and theta.
