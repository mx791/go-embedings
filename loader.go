package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

var old_to_new_indexes = make(map[int]int)
var new_to_old_indexes = make(map[int]int)
var freqs = make(map[int]int)

var ITERS = 13
var LEARN_RATE = 0.1
var LEARN_RATE_DECAY = 0.3
var EMBEDING_SIZE = 48
var BATCH_SIZE = 5
var MAX_ID = 3000000
var FILE_PATH = "./playlists_.json"
var OUTPUT_PATH = "./embedings.json"

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
				freqs[list_value] = 1
				counter += 1
			} else {
				freqs[list_value] = freqs[list_value] + 1
			}
		}
	}
	fmt.Println(len(old_to_new_indexes), len(new_to_old_indexes))
	return new_liste
}

// fonction valeur absolue
func Abs(number float64) float64 {
	if number > 0 {
		return number
	} else {
		return -number
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

// selectionne un id aléatoire depuis la même liste qu'un autre id
func select_true_item(liste []int, item_id int) int {
	selected_id := liste[rand.Intn(len(liste))]
	for selected_id == item_id || selected_id >= MAX_ID {
		selected_id = liste[rand.Intn(len(liste))]
	}
	return selected_id
}

// sélectionne un id aléatoire parmis tous les id dispo
func gen_id(listes [][]int) int {
	vect_id := rand.Intn(len(listes))
	new_item_id := listes[vect_id][rand.Intn(len(listes[vect_id]))]
	for new_item_id >= MAX_ID {
		vect_id = rand.Intn(len(listes))
		new_item_id = listes[vect_id][rand.Intn(len(listes[vect_id]))]
	}
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
	for value_id := 0; value_id < EMBEDING_SIZE; value_id++ {

		if anchor[value_id] > similar[value_id] {
			similar[value_id] += LEARN_RATE
		} else {
			similar[value_id] -= LEARN_RATE
		}

		if anchor[value_id] > different[value_id] {
			different[value_id] -= LEARN_RATE
		} else {
			different[value_id] += LEARN_RATE
		}
	}
}

// serieusement besoin d'expliquer ?
func train_embedings(liste [][]int) [][]float64 {
	embeding_matrix := init_embedings(MAX_ID, EMBEDING_SIZE)
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

					base_vector := embeding_matrix[anchor_id]
					similar_vector := embeding_matrix[similar_item_id]
					different_vector := embeding_matrix[different_item_id]

					inside_dst += euclidian_dst(base_vector, similar_vector)
					inside_sim += cosin_sim(base_vector, similar_vector)

					outside_dst += euclidian_dst(base_vector, different_vector)
					outside_sim += cosin_sim(base_vector, different_vector)

					update_embeding(base_vector, similar_vector, different_vector)
				}
				wg.Done()
			}(list_id)
		}

		wg.Wait()

		for embeding_id := 0; embeding_id < len(embeding_matrix); embeding_id++ {
			normalize_vect(embeding_matrix[embeding_id])
		}
		LEARN_RATE = LEARN_RATE * LEARN_RATE_DECAY
		fmt.Println("iter", iter+1, "/", ITERS, "cosin_sim=", (outside_sim / inside_sim), "euclidian_dst=", (inside_dst / outside_dst), "learning_rate=", LEARN_RATE)
	}
	return embeding_matrix
}

// normalize un vecteur pour que sa norme soit 1
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

func create_random_adjacency_list(size int, max_index int, max_vect_size int, min_vect_size int) [][]int {
	return_list := make([][]int, size)
	for line_id := range return_list {
		return_list[line_id] = make([]int, rand.Intn(max_vect_size-min_vect_size)+min_vect_size)
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

func load_liste_file(path string) [][]int {
	jsonFile, err := os.Open(path)
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var datas [][]int
	json.Unmarshal(byteValue, &datas)
	return datas
}

func find_closest(vector []float64, n int, matrix [][]float64) {
	cc := 0
	stepp := 0.0002
	for seuil := 0.0; seuil < 1.0; seuil += stepp {
		for item := 0; item < MAX_ID; item++ {
			dst := euclidian_dst(vector, matrix[item])
			if dst > seuil && dst < seuil+stepp {
				fmt.Println(dst, item)
				cc += 1
			}
			if cc > n {
				return
			}
		}
	}
}

func clear_liste(liste [][]int) [][]int {
	filtered_liste := make([][]int, 0)
	for i := 0; i < len(liste); i++ {
		vect := make([]int, 0)
		used_indexes := make(map[int]int)
		for e := 0; e < len(liste[i]); e++ {
			value := liste[i][e]
			_, ok := used_indexes[value]
			if value < MAX_ID && !ok {
				vect = append(vect, value)
				used_indexes[value] = 1
			}
		}
		if len(vect) > 3 {
			filtered_liste = append(filtered_liste, vect)
		}
	}
	return filtered_liste
}

// renvoie les id min/max de toutes les listes
func getMinMaxId(liste [][]int) (int, int) {
	min := liste[0][0]
	max := liste[0][0]
	for i := 0; i < len(liste); i++ {
		for e := 0; e < len(liste[i]); e++ {
			if liste[i][e] > max {
				max = liste[i][e]
			} else if liste[i][e] < min {
				min = liste[i][e]
			}
		}
	}
	return min, max
}

func parse_param() {
	argsWithProg := os.Args
	for i := 1; i < len(argsWithProg); i++ {
		arg_array := strings.Split(argsWithProg[i], "=")
		if len(arg_array) == 1 {
			continue
		}
		value, err := strconv.Atoi(arg_array[1])
		if err != nil {
			if arg_array[0] == "FILE_PATH" {
				fmt.Println("seting FILE_PATH to", arg_array[1])
				FILE_PATH = arg_array[1]
			} else if arg_array[0] == "OUTPUT_PATH" {
				fmt.Println("seting OUTPUT_PATH to", arg_array[1])
				OUTPUT_PATH = arg_array[1]
			}
			continue
		}
		if arg_array[0] == "ITERS" {
			ITERS = value
			fmt.Println("seting ITERS to", arg_array[1])
		} else if arg_array[0] == "EMBEDING_SIZE" {
			EMBEDING_SIZE = value
			fmt.Println("seting EMBEDING_SIZE to", arg_array[1])
		} else if arg_array[0] == "MAX_ID" {
			MAX_ID = value
			fmt.Println("seting MAX_ID to", arg_array[1])
		}
	}
}

func save_datas(matrix [][]float64, target_path string) {
	file, _ := json.MarshalIndent(matrix, "", " ")
	_ = ioutil.WriteFile(target_path, file, 0644)
}

func main() {

	parse_param()

	datas := load_liste_file(FILE_PATH)
	fmt.Println(len(datas), "listes")
	datas = clear_liste(datas)
	fmt.Println(len(datas), "listes après filtrage")

	min, max := getMinMaxId(datas)
	fmt.Println("min id:", min, "max id:", max)
	MAX_ID = max + 1

	// build_normalized_list(datas)
	matrix := train_embedings(datas)
	save_datas(matrix, OUTPUT_PATH)
}
