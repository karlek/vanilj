package mandel

// isMemberOfSet determines if the complex point z is member of the mandelbrot
// set.
func isMemberOfSet(z complex128) bool {
	// same as 2 > cmplx.Abs(z)
	return real(z)*real(z)+imag(z)*imag(z) < 4
}

// isInBulb is a computational saving function for determining if points are
// inside one of the main mandelbrot bulbs.
// Credits: https://github.com/morcmarc/buddhabrot/blob/master/buddhabrot.go
func isInBulb(c complex128) bool {
	Cr, Ci := real(c), imag(c)
	// Main cardioid
	if !(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))*(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))+(Cr-0.25)) < 0.25*Ci*Ci) {
		// 2nd order period bulb
		if !((Cr+1.0)*(Cr+1.0)+(Ci*Ci) < 0.0625) {
			// smaller bulb left of the period-2 bulb
			if !((((Cr + 1.309) * (Cr + 1.309)) + Ci*Ci) < 0.00345) {
				// smaller bulb bottom of the main cardioid
				if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci-0.744)*(Ci-0.744)) < 0.0088) {
					// smaller bulb top of the main cardioid
					if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci+0.744)*(Ci+0.744)) < 0.0088) {
						return false
					}
				}
			}
		}
	}
	return true
}

// divergence returns the number of iterations it takes for a complex point
// to leave the mandelbrot set and also returns the point last point (which
// could be outside the mandelbrot set).
func divergence(c complex128, iterations int64) (i int64, z complex128) {
	// Ignore points in the main bulbs (the never diverge).
	if isInBulb(c) {
		return int64(iterations), c
	}

	// Saved value for cycle-detection.
	var bfract complex128

	z = complex(0, 0)
	for i = 0; i < int64(iterations); i += 1 {
		z = z*z + c
		if !isMemberOfSet(z) {
			break
		}
		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			bfract = z
		} else if z == bfract {
			return int64(iterations), z
		}
	}
	return i, z
}
