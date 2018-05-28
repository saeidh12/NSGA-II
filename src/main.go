package main

import (
	"fmt"
	"math/rand"
	"time"

	"./nsga2"
)

const number_of_repetition = 200
const population_size = 100
const number_of_neihborhoods = 50
const number_of_firestation_locations = 40
const init_p_mutation = 0.3
const init_p_crossover = 0.7
const x_axis_count = 20
const y_axis_count = 50

var list_of_points [][][]int

func stopCondition(generation int, children [][]int) bool {
	fmt.Println("gen: ", generation, "answer vec len: ", nsga2.EuclideanDist(children[0], []int{0, 0}))
	if generation == number_of_repetition {
		return true
	} else {
		return false
	}
}

func ObjectFunction1(slice []int) float64 {
	number_of_active_stations := 0
	changed_station := false
	for _, item := range slice {
		if item < number_of_neihborhoods {
			changed_station = false
		} else {
			if !changed_station {
				number_of_active_stations++
			}
			changed_station = true
		}
	}
	if !changed_station {
		number_of_active_stations++
	}
	return float64(number_of_active_stations)
}

func ObjectFunction2(slice []int) float64 {
	neibors := 0
	sum := 0.0
	for index, item := range slice {
		if item < number_of_neihborhoods {
			neibors++
		} else {
			for i := index - neibors; i < index; i++ {
				hood_coord := list_of_points[0][slice[i]]
				firestation := list_of_points[1][item-number_of_neihborhoods]
				sum += nsga2.EuclideanDist(hood_coord, firestation)
			}
			neibors = 0
		}
	}
	index := len(slice)
	for i := index - neibors; i < index; i++ {
		hood_coord := list_of_points[0][slice[i]]
		firestation := list_of_points[1][number_of_firestation_locations-1]
		sum += nsga2.EuclideanDist(hood_coord, firestation)
	}

	return sum
}

func InitSpace() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	chosen_points := [][]int{}
	list_of_points := [][][]int{}
	list_of_points = append(list_of_points, [][]int{})
	list_of_points = append(list_of_points, [][]int{})

	for i := 0; i < number_of_neihborhoods; {
		rand_x := r1.Intn(x_axis_count)
		rand_y := r1.Intn(y_axis_count)
		new_neigborhood := []int{rand_x, rand_y}

		if nsga2.SliceIndexInMat(new_neigborhood, chosen_points) == -1 {
			list_of_points[0] = append(list_of_points[0], new_neigborhood)
			chosen_points = append(chosen_points, new_neigborhood)
			i++
		}
	}

	for i := 0; i < number_of_firestation_locations; {
		rand_x := r1.Intn(x_axis_count)
		rand_y := r1.Intn(y_axis_count)
		new_firestation := []int{rand_x, rand_y}

		if nsga2.SliceIndexInMat(new_firestation, chosen_points) == -1 {
			list_of_points[1] = append(list_of_points[1], new_firestation)
			chosen_points = append(chosen_points, new_firestation)
			i++
		}
	}
}

func main() {
	InitSpace()
	object_funcs := []func([]int) float64{ObjectFunction1, ObjectFunction2}
	problem_size := number_of_neihborhoods + number_of_firestation_locations - 1
	results := nsga2.NSGA2(
		stopCondition,
		object_funcs,
		population_size,
		problem_size,
		init_p_crossover,
		init_p_mutation)
	fmt.Println("answer vec: ", results[0])
	fmt.Println("answer vec len: ", nsga2.EuclideanDist(results[0], []int{0, 0}))
}
