package main

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

var points [iterations]complex128

func divergencePrim(c complex128) []complex128 {
	if isInBulb(c) {
		return nil
	}

	var brent complex128
	z := complex(0, 0)
	var num int
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Cycle detection.
		if (i-1)&i == 0 && i > 1 {
			brent = z
		} else if z == brent {
			return nil
		}
		// Diverges.
		if x, y := real(z), imag(z); x*x+y*y >= bailout {
			return points[:num]
		}
		points[num] = z
		num++
	}
	// Converges.
	return nil
}
