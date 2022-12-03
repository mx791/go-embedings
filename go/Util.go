package main

// vector normalisation
func normalize_vect(vect []float64) {
	norm := 0.0
	for i := 0; i < len(vect); i++ {
		norm += Abs(vect[i])
	}
	if norm == 0.0 {
		norm = 1.0
	}
	for i := 0; i < len(vect); i++ {
		vect[i] /= norm
	}
}

func Abs(number float64) float64 {
	if number > 0 {
		return number
	} else {
		return -number
	}
}
