package mandel

// isMemberOfSet determines if the complex point z is member of the mandelbrot
// set.
func isMemberOfSet(z complex128) bool {
	// same as 2 > cmplx.Abs(z)
	return real(z)*real(z)+imag(z)*imag(z) < 4
}

// isInBulb is a computational saving function for determining if points are
// inside one of the main mandelbrot bulbs.
func isInBulb(c complex128) bool {
	x, y := real(c), imag(c)
	q := (x-1.0/4.0)*(x-1.0/4.0) + y*y
	return q*(q+(x-1.0/4.0)) < 1.0/4.0*y*y
}

// divergence returns the number of iterations it takes for a complex point
// to leave the mandelbrot set and also returns the point last point (which
// could be outside the mandelbrot set).
func divergence(c complex128, iterations float64) (i float64, z complex128) {
	// Ignore points in the main bulbs (the never diverge).
	if isInBulb(c) {
		return iterations, c
	}

	z = complex(0, 0)
	for i = 0.0; i < iterations; i += 1 {
		z = z*z + c
		if !isMemberOfSet(z) {
			break
		}
	}
	return i, z
}
