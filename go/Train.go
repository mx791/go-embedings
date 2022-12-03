package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Initialise embeddings with random values
func init_embeddings(count int, n_dim int) [][]float64 {
	embedding_matrix := make([][]float64, count)
	for idx := 0; idx < count; idx++ {
		embedding_matrix[idx] = make([]float64, n_dim)
		for idx2 := 0; idx2 < n_dim; idx2++ {
			embedding_matrix[idx][idx2] = 0.5 - rand.Float64()
		}
	}
	return embedding_matrix
}

// select a track from the same playlist than another
func select_true_item(liste []int, item_id int) int {
	selected_id := liste[rand.Intn(len(liste))]
	for selected_id == item_id || selected_id >= MAX_ID {
		selected_id = liste[rand.Intn(len(liste))]
	}
	return selected_id
}

// select randomly an id
func gen_id(listes [][]int) int {
	vect_id := rand.Intn(len(listes))
	new_item_id := listes[vect_id][rand.Intn(len(listes[vect_id]))]
	for new_item_id >= MAX_ID {
		vect_id = rand.Intn(len(listes))
		new_item_id = listes[vect_id][rand.Intn(len(listes[vect_id]))]
	}
	return new_item_id
}

// search if a value is in an array
func in_array(liste []int, value int) bool {
	for array_value := range liste {
		if value == array_value {
			return true
		}
	}
	return false
}

// select an id which is NOT in a given list
func select_false_item(listes [][]int, list_id int, item_id int) int {
	new_item_id := gen_id(listes)
	for new_item_id == item_id || in_array(listes[list_id], new_item_id) {
		new_item_id = gen_id(listes)
	}
	return new_item_id
}

// Update an ambedding based on a track, a similar track and a different track
func update_embedding(anchor []float64, similar []float64, different []float64) {
	for value_id := 0; value_id < EMBEDDING_SIZE; value_id++ {

		similar[value_id] += (anchor[value_id] - similar[value_id]) * LEARN_RATE

		if anchor[value_id] > different[value_id] {
			different[value_id] -= LEARN_RATE
		} else {
			different[value_id] += LEARN_RATE
		}
	}
}

// main loop for embedding training
func train_embeddings(liste [][]int) [][]float64 {
	embedding_matrix := init_embeddings(MAX_ID, EMBEDDING_SIZE)
	for iter := 0; iter < ITERS; iter++ {

		inside_dst := 0.0
		inside_sim := 0.0
		outside_dst := 0.0
		outside_sim := 0.0

		var wg sync.WaitGroup

		for list_id := 0; list_id < len(liste); list_id++ {
			wg.Add(1)
			go func(list_id int) {
				for item_id := 0; item_id < len(liste[list_id]); item_id++ {

					anchor_id := liste[list_id][item_id]

					similar_item_id := select_true_item(liste[list_id], anchor_id)
					different_item_id := select_false_item(liste, list_id, anchor_id)

					base_vector := embedding_matrix[anchor_id]
					similar_vector := embedding_matrix[similar_item_id]
					different_vector := embedding_matrix[different_item_id]

					inside_dst += euclidian_dst(base_vector, similar_vector)
					inside_sim += cosin_sim(base_vector, similar_vector)

					outside_dst += euclidian_dst(base_vector, different_vector)
					outside_sim += cosin_sim(base_vector, different_vector)

					update_embedding(base_vector, similar_vector, different_vector)
				}
				wg.Done()
			}(list_id)
		}

		wg.Wait()

		LEARN_RATE = LEARN_RATE * LEARN_RATE_DECAY
		fmt.Println("iter", iter+1, "/", ITERS, "cosin_sim=", (outside_sim / inside_sim), "euclidian_dst=", (inside_dst / outside_dst), "learning_rate=", LEARN_RATE)
	}
	return embedding_matrix
}
