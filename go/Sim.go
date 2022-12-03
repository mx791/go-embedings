package main

// fonction gnérique de calcul de similarité entre deux vecteurs
func sim(vect_a []float64, vect_b []float64, product bool) float64 {
	norm_a := 0.0
	norm_b := 0.0
	sim := 0.0

	for i := 0; i < len(vect_b); i++ {
		if product {
			sim += Abs(vect_a[i] * vect_b[i])
		} else {
			sim += (vect_a[i] - vect_b[i]) * (vect_a[i] - vect_b[i])
		}
		norm_a += Abs(vect_a[i])
		norm_b += Abs(vect_b[i])
	}

	if norm_a == 0.0 {
		norm_a = 1.0
	}
	if norm_b == 0.0 {
		norm_b = 1.0
	}

	return sim / (norm_a * norm_b)
}

// similarité cosinus, surcouche de sim
func cosin_sim(vect_a []float64, vect_b []float64) float64 {
	return sim(vect_a, vect_b, true)
}

// distance euclidienne cosinus, surcouche de sim
func euclidian_dst(vect_a []float64, vect_b []float64) float64 {
	return sim(vect_a, vect_b, false)
}
