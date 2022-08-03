package main;


import (
	"fmt"
	"math/rand"
)

var old_to_new_indexes = make(map[int]int)
var new_to_old_indexes = make(map[int]int)
var freqs = make(map[int]int)

var ITERS = 20
var LEARN_RATE = 0.01
var EMBEDING_SIZE = 4
var BATCH_SIZE = 5

func build_normalized_list(liste [][]int) [][]int {
	new_liste := make([][]int, len(liste))
	counter := 0
	for list_id := 0; list_id < len(liste); list_id++ {
		for value_id := 0; value_id < len(liste[list_id]); value_id++ {
			list_value := liste[list_id][value_id]
			_, ok := freqs[list_value]
			if !ok {
				old_to_new_indexes[list_value] = counter
				new_to_old_indexes[counter] = list_value
				freqs[counter] = 1
				counter += 1
			} else {
				freqs[list_value] += 1
			}
		}
	}
	return new_liste
}

// fonction valeur absolue
func Abs(number float64) float64 {
	if number > 0 {
		return number
	} else {
		return - number
	}
}

// créer des embedings aléatoires
func init_embedings(count int, n_dim int) [][]float64 {
	embeding_matrix := make([][]float64, count)
	for idx := 0; idx < count; idx++ {
		embeding_matrix[idx] = make([]float64, n_dim)
		for idx2 := 0; idx2 < n_dim; idx2++ {
			embeding_matrix[idx][idx2] = 0.5 - rand.Float64()
		}
	}
	return embeding_matrix
}

// fonction gnérique de calcul de similarité entre deux vecteurs
func sim(vect_a []float64, vect_b []float64, product bool) float64 {
	norm_a := 0.0
	norm_b := 0.0
	sim := 0.0

	for i:=0; i<len(vect_b); i++ {
		if product {
			sim += Abs(vect_a[i] * vect_b[i])
		} else {
			sim += Abs(vect_a[i] - vect_b[i])
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

// selectionne un id aléatoire depuis la même liste qu'un autre id
func select_true_item(liste []int, item_id int) int {
	selected_id := liste[rand.Intn(len(liste))]
	for selected_id == item_id {
		selected_id = liste[rand.Intn(len(liste))]
	}
	return selected_id
}

// sélectionne un id aléatoire parmis tous les id dispo
func gen_id(listes [][]int) int {
	vect_id := rand.Intn(len(listes))
	new_item_id := listes[vect_id][rand.Intn(len(listes[vect_id]))]
	return new_item_id
}

// recherche de présence dans un tableau
func in_array(liste []int, value int) bool {
	for array_value := range liste {
		if value == array_value {
			return true
		}
	}
	return false
}

// selectionne un id absent d'une liste donnée
func select_false_item(listes [][]int, list_id int, item_id int) int {
	new_item_id := gen_id(listes)
	for new_item_id == item_id || in_array(listes[list_id], new_item_id) {
		new_item_id = gen_id(listes)
	}
	return new_item_id
}

// met à jour un embeding pour l'eloigner des items différents et le rapprocher des similaires
func update_embeding(anchor []float64, similar []float64, different []float64) {
	for value_id:=0; value_id<EMBEDING_SIZE; value_id++ {
		if anchor[value_id] > similar[value_id] {
			anchor[value_id] -= LEARN_RATE
		} else {
			anchor[value_id] += LEARN_RATE
		}

		if anchor[value_id] > different[value_id] {
			different[value_id] -= LEARN_RATE
		} else {
			different[value_id] += LEARN_RATE
		}
	}
}

func train_embedings(liste [][]int) [][]float64 {
	embeding_matrix := init_embedings(len(new_to_old_indexes), EMBEDING_SIZE)
	for iter := 0; iter < ITERS; iter++ {

		inside_dst := 0.0
		inside_sim := 0.0
		outside_dst := 0.0
		outside_sim := 0.0

		for list_id := 0; list_id < len(liste); list_id++ {

			for item_id:=0; item_id < len(liste[list_id]); item_id++ {

				similar_item_id := select_true_item(liste[list_id], item_id)
				different_item_id := select_false_item(liste, list_id, item_id)
	
				base_vector := embeding_matrix[item_id]
				similar_vector := embeding_matrix[similar_item_id]
				different_vector := embeding_matrix[different_item_id]
	
				inside_dst += euclidian_dst(base_vector, similar_vector)
				inside_sim += cosin_sim(base_vector, similar_vector)
	
				outside_dst += euclidian_dst(base_vector, different_vector)
				outside_sim += cosin_sim(base_vector, different_vector)
				
				update_embeding(base_vector, similar_vector, different_vector)
			}
		}

		fmt.Println((iter, "/", ITERS, ", cosin_sim=", (outside_sim/inside_sim), ", euclidian_dst=", (inside_dst/outside_dst))
	}
	return embeding_matrix
}

// normalize un vecteur pour que sa norme soit 1
func normalize_vect(vect []float64) {
	norm := 0.0
	for i:=0; i<len(vect); i++ {
		norm += Abs(vect[i])
	}
	if norm == 0.0 {
		norm = 1.0
	}
	for i:=0; i<len(vect); i++ {
		vect[i] /= norm
	}
}

func create_random_adjacency_list(size int, max_index int, max_vect_size int, min_vect_size int) [][]int {
	return_list := make([][]int, size)
	for line_id := range return_list {
		return_list[line_id] = make([]int, rand.Intn(max_vect_size-min_vect_size) + min_vect_size)
		used_indexes := make(map[int]int)
		
		for pos_id := range return_list[line_id] {
			id := rand.Intn(max_index)
			_, ok := used_indexes[id]
			for ok {
				id = rand.Intn(max_index)
				_, ok = used_indexes[id]
			}
			used_indexes[id] = 1
			return_list[line_id][pos_id] = id
		}
	}
	return return_list
}

func main() {
	test_list := create_random_adjacency_list(15000, 50, 10, 3)
	build_normalized_list(test_list)
	train_embedings(test_list)
}
