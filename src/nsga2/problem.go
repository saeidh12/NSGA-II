package nsga2

import (
	"fmt"
	"math/rand"
)

type Problem struct {
	list_of_points                  [][][]int
	number_of_repetition            int
	population_size                 int
	problem_size                    int
	number_of_neihborhoods          int
	number_of_firestation_locations int
	init_p_selected                 float64
	init_p_mutation                 float64
	init_p_crossover                float64
	x_axis_count                    int
	y_axis_count                    int
}

func (p *Problem) Init(
	number_of_repetition,
	population_size,
	number_of_neihborhoods,
	number_of_firestation_locations,
	x_axis_count,
	y_axis_count int,
	init_p_mutation,
	init_p_crossover,
	init_p_selected float64,
) {
	p.number_of_repetition = number_of_repetition
	p.population_size = population_size
	p.number_of_neihborhoods = number_of_neihborhoods
	p.number_of_firestation_locations = number_of_firestation_locations
	p.x_axis_count = x_axis_count
	p.y_axis_count = y_axis_count
	p.init_p_mutation = init_p_mutation
	p.init_p_crossover = init_p_crossover
	p.init_p_selected = init_p_selected
	p.problem_size = number_of_neihborhoods + number_of_firestation_locations - 1

	// s1 := rand.NewSource(time.Now().UnixNano())
	// r1 := rand.New(s1)

	chosen_points := [][]int{}
	p.list_of_points = append(p.list_of_points, [][]int{})
	p.list_of_points = append(p.list_of_points, [][]int{})

	for i := 0; i < p.number_of_neihborhoods; {
		rand_x := rand.Intn(p.x_axis_count) // r1
		rand_y := rand.Intn(p.y_axis_count) // r1
		new_neigborhood := []int{rand_x, rand_y}

		if SliceIndexInMat(new_neigborhood, chosen_points) == -1 {
			p.list_of_points[0] = append(p.list_of_points[0], new_neigborhood)
			chosen_points = append(chosen_points, new_neigborhood)
			i++
		}
	}

	for i := 0; i < p.number_of_firestation_locations; {
		rand_x := rand.Intn(p.x_axis_count) // r1
		rand_y := rand.Intn(p.y_axis_count) // r1
		new_firestation := []int{rand_x, rand_y}

		if SliceIndexInMat(new_firestation, chosen_points) == -1 {
			p.list_of_points[1] = append(p.list_of_points[1], new_firestation)
			chosen_points = append(chosen_points, new_firestation)
			i++
		}
	}
}

func (p Problem) stopCondition(generation int, children [][]int, mapped_mat [][]float64) bool {
	fmt.Println("gen: ", generation, " number of stations: ", mapped_mat[0][0]/100, "  avg vec len: ", AvgEuclideanDist(mapped_mat), "  answer vec len: ", EuclideanDistFloat(mapped_mat[0], []float64{0, 0}))

	Plot(generation, children, mapped_mat)

	if generation == p.number_of_repetition {
		return true
	} else {
		return false
	}
}

func (p Problem) ObjectFunction1(slice []int) float64 {
	number_of_active_stations := 0
	changed_station := false
	for _, item := range slice {
		if item < p.number_of_neihborhoods {
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
	return float64(number_of_active_stations * 100)
}

func (p Problem) ObjectFunction2(slice []int) float64 {
	neibors := 0
	sum := 0.0
	for index, item := range slice {
		if item < p.number_of_neihborhoods {
			neibors++
		} else {
			for i := index - neibors; i < index; i++ {
				hood_coord := p.list_of_points[0][slice[i]]
				firestation := p.list_of_points[1][item-p.number_of_neihborhoods]
				sum += EuclideanDist(hood_coord, firestation)
			}
			neibors = 0
		}
	}
	index := len(slice)
	for i := index - neibors; i < index; i++ {
		hood_coord := p.list_of_points[0][slice[i]]
		firestation := p.list_of_points[1][p.number_of_firestation_locations-1]
		sum += EuclideanDist(hood_coord, firestation)
	}

	return sum
}
