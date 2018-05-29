package main

import (
	"fmt"

	"./nsga2"
)

const number_of_repetition = 500
const population_size = 10
const number_of_neihborhoods = 50          // 50
const number_of_firestation_locations = 20 // 40
const init_p_selected = 0.7
const init_p_mutation = 0.7
const init_p_crossover = 0.7
const x_axis_count = 20
const y_axis_count = 50

func main() {
	problem := nsga2.Problem{}
	problem.Init(
		number_of_repetition,
		population_size,
		number_of_neihborhoods,
		number_of_firestation_locations,
		x_axis_count,
		y_axis_count,
		init_p_mutation,
		init_p_crossover,
		init_p_selected,
	)
	results := nsga2.NSGA2(problem)
	// fmt.Println("answer vec: ", results[0])
	fmt.Println("answer vec len: ", nsga2.EuclideanDist(results[0], []int{0, 0}))
	fmt.Println("avg vec len: ", nsga2.AvgEuclideanDist(results))
}
