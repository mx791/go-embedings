package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var old_to_new_indexes = make(map[int]int)
var new_to_old_indexes = make(map[int]int)
var freqs = make(map[int]int)

var ITERS = 50
var LEARN_RATE = 0.1
var LEARN_RATE_DECAY = 0.3
var EMBEDDING_SIZE = 32
var MAX_ID = 3000000
var FILE_PATH = "../output/playlists_.json"
var OUTPUT_PATH = "../output/embedings.json"

// normalise all ids, so they start with id 0 and and with mAX_ID, not used here as we have clean data
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

// read the JSON file
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

// pre-process the input list of list
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
		if len(vect) >= 3 {
			filtered_liste = append(filtered_liste, vect)
		}
	}
	return filtered_liste
}

// find min and max ids from all lists
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

// parse command line args
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
			EMBEDDING_SIZE = value
			fmt.Println("seting EMBEDING_SIZE to", arg_array[1])
		} else if arg_array[0] == "MAX_ID" {
			MAX_ID = value
			fmt.Println("seting MAX_ID to", arg_array[1])
		}
	}
}

// save embedding in a json file
func save_datas(matrix [][]float64, target_path string) {
	file, _ := json.MarshalIndent(matrix, "", " ")
	_ = ioutil.WriteFile(target_path, file, 0644)
}

func main() {

	parse_param()

	datas := load_liste_file(FILE_PATH)
	fmt.Println(len(datas), "lists")
	datas = clear_liste(datas)
	fmt.Println(len(datas), "lists after filtering")

	min, max := getMinMaxId(datas)
	fmt.Println("min id:", min, "max id:", max)
	MAX_ID = max + 1

	matrix := train_embeddings(datas)
	save_datas(matrix, OUTPUT_PATH)
}
