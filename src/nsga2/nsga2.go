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
			if SliceLT(item, mapped_mat[i+1]) {
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
		if SliceLT(mapped_mat[i], mapped_mat[i-1]) {
			result[index] = append(result[index], mat[i])
		} else {
			index++
			result = append(result, [][]int{})
			result[index] = append(result[index], mat[i])
		}
	}

	return result
}

func SelectParentsByRank(population [][]int, population_size int, p_selected float64) [][]int {
	num_selected := int(p_selected * float64(population_size))
	return population[:num_selected]
}

func CrossoverAndMutation(selected [][]int, p_crossover, p_mutation float64) [][]int {
	result := [][]int{}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result = append(result, selected...)

	for _, item := range selected[:len(selected)-1] {
		random := r1.Float64()
		if random < p_crossover {
			// child1, child2 := Crossover1(item, selected[index + 1])
			// result = append(result, child1, child2)

			child := Crossover(item)
			if SliceIndexInMat(child, result) == -1 {
				result = append(result, child)
			}
		}

		if random > (1 - p_mutation) {
			child := Mutate(item)
			if SliceIndexInMat(child, result) == -1 {
				result = append(result, child)
			}
		}
	}
	return result
}

func NSGA2(problem Problem) [][]int {
	obj_func := []func([]int) float64{problem.ObjectFunction1, problem.ObjectFunction2}
	// fmt.Println("here 0")
	population := InitializePopulation(problem.population_size, problem.problem_size)
	// fmt.Println("here 1")
	mapped_mat := MapMatWithObjectiveFunc(population, obj_func)
	// fmt.Println("here 2")
	population = SliceFastNondominatedSort(population, mapped_mat)
	// fmt.Println("here 3")
	selected := SelectParentsByRank(population, problem.population_size, problem.init_p_selected)
	// fmt.Println("here 4")
	children := CrossoverAndMutation(selected, problem.init_p_crossover, problem.init_p_mutation)
	// fmt.Println("here 5")

	for generation := 0; !problem.stopCondition(generation, children); generation++ {
		// union := append(population, children...)
		union := [][]int{}
		for _, item := range append(population, children...) {
			if SliceIndexInMat(item, union) == -1 {
				item_copy := make([]int, len(item))
				copy(item_copy, item)
				union = append(union, item_copy)
			}
		}
		// fmt.Println("here 6")
		mapped_mat = MapMatWithObjectiveFunc(union, obj_func)
		// fmt.Println("here 7")
		union = SliceFastNondominatedSort(union, mapped_mat)
		// fmt.Println("here 8")
		fronts := CreateFronts(union, mapped_mat)
		// fmt.Println("here 9")

		parents := [][]int{}
		front_l := -1
		// fmt.Println("here 10")
		for index, front := range fronts {
			// fmt.Println("here 11")
			if len(parents)+len(front) > problem.population_size {
				// fmt.Println("here 12")
				front_l = index
				break
			} else {
				// fmt.Println("here 13")
				parents = append(parents, front...)
			}
		}

		amount_left := problem.population_size - len(parents)
		if amount_left != 0 {
			// fmt.Println("here 14", front_l, amount_left, fronts)
			parents = append(parents, fronts[front_l][:amount_left]...)
		}
		// fmt.Println("here 15")
		selected = SelectParentsByRank(parents, problem.population_size, problem.init_p_selected)
		// fmt.Println("here 16")
		population = children
		// fmt.Println("here 17")
		children = CrossoverAndMutation(selected, problem.init_p_crossover, problem.init_p_mutation)
		// fmt.Println("here 18")
	}

	mapped_mat = MapMatWithObjectiveFunc(children, obj_func)
	// fmt.Println("here 19")
	children = SliceFastNondominatedSort(children, mapped_mat)
	// fmt.Println("here 20")
	return children
}
