package main

import (
	"fmt"

	"./nsga2"
)

const number_of_repetition = 100
const population_size = 200
const number_of_neihborhoods = 70
const number_of_firestation_locations = 70
const init_p_selected = 0.3
const init_p_mutation = 0.2
const init_p_crossover = 0.4
const x_axis_count = 100
const y_axis_count = 100

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
	result := nsga2.NSGA2(problem)
	fmt.Println("number of stations: ", problem.ObjectFunction1(result)/100, " answer: ", result)
}
