package nsga2

import (
	"math/rand"
	"time"
)

func InitializePopulation(population_size, problem_size int) [][]int {
	population := [][]int{}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < population_size; {
		rand_slice := r1.Perm(problem_size)

		if SliceIndexInMat(rand_slice, population) == -1 {
			population = append(population, rand_slice)
			i++
		}
	}

	return population
}

func MapMatWithObjectiveFunc(mat [][]int, obj_funcs []func([]int) float64) [][]float64 {
	func_on_mat := [][]float64{}

	for _, slice := range mat {
		mapped_slice := FuncsOnSlice(slice, obj_funcs)
		func_on_mat = append(func_on_mat, mapped_slice)
	}

	return func_on_mat
}

func SliceFastNondominatedSort(mat [][]int, mapped_mat [][]float64) [][]int {
	swaps := 1
	for swaps != 0 {
		swaps = 0
		for i, item := range mapped_mat[:len(mapped_mat)-1] {
			if SliceLTE(item, mapped_mat[i+1]) {
				mat[i], mat[i+1] = mat[i+1], mat[i]
				mapped_mat[i], mapped_mat[i+1] = mapped_mat[i+1], mapped_mat[i]
				swaps++
			}
		}
	}

	return mat
}

func CreateFronts(mat [][]int, mapped_mat [][]float64) [][][]int {
	index := 0
	result := [][][]int{}
	result = append(result, [][]int{})
	result[index] = append(result[index], mat[0])

	for i := 1; i < len(mat); i++ {
		if SliceLTE(mapped_mat[i], mapped_mat[i-1]) {
			result[index] = append(result[index], mat[i])
		} else {
			index++
			result = append(result, [][]int{})
			result[index] = append(result[index], mat[i])
		}
	}

	return result
}

func SelectParentsByRank(population [][]int, population_size int) [][]int {
	p_selected := 0.3
	num_selected := int(p_selected * float64(population_size))
	return population[:num_selected]
}

func CrossoverAndMutation(selected [][]int, p_crossover, p_mutation float64) [][]int {
	result := [][]int{}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for _, item := range selected[:len(selected)-1] {
		random := r1.Float64()
		if random < p_crossover {
			// child1, child2 := Crossover1(item, selected[index + 1])
			// result = append(result, child1, child2)

			child := Crossover(item)
			result = append(result, child)
		}

		if random > (1 - p_mutation) {
			child := Mutate(item)
			result = append(result, child)
		}
	}
	return result
}

func NSGA2(stopCondition func(int, [][]int) bool,
	obj_funcs []func([]int) float64,
	population_size, problem_size int,
	init_p_crossover, init_p_mutation float64) [][]int {
	population := InitializePopulation(population_size, problem_size)
	mapped_mat := MapMatWithObjectiveFunc(population, obj_funcs)
	population = SliceFastNondominatedSort(population, mapped_mat)
	selected := SelectParentsByRank(population, population_size)
	children := CrossoverAndMutation(selected, init_p_crossover, init_p_mutation)

	for generation := 0; stopCondition(generation, children); generation++ {
		union := append(population, children...)
		mapped_mat = MapMatWithObjectiveFunc(union, obj_funcs)
		union = SliceFastNondominatedSort(union, mapped_mat)
		fronts := CreateFronts(union, mapped_mat)

		parents := [][]int{}
		front_l := -1

		for index, front := range fronts {
			if len(parents)+len(front) > population_size {
				front_l = index
				break
			} else {
				parents = append(parents, front...)
			}
		}

		amount_left := population_size - len(parents)
		if amount_left != 0 {
			parents = append(parents, fronts[front_l][:amount_left]...)
		}

		selected = SelectParentsByRank(parents, population_size)
		population = children
		children = CrossoverAndMutation(selected, init_p_crossover, init_p_mutation)
	}

	mapped_mat = MapMatWithObjectiveFunc(children, obj_funcs)
	children = SliceFastNondominatedSort(children, mapped_mat)
	return children
}
