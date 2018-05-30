package nsga2

import (
	"fmt"
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

func SliceFastNondominatedSort(mat *[][]int, mapped_mat *[][]float64) {
	swaps := 1
	for swaps != 0 {
		swaps = 0
		for i := 0; i < len(*mapped_mat)-1; i++ {
			if EuclideanLT((*mapped_mat)[i+1], (*mapped_mat)[i]) {
				(*mat)[i], (*mat)[i+1] = (*mat)[i+1], (*mat)[i]
				(*mapped_mat)[i], (*mapped_mat)[i+1] = (*mapped_mat)[i+1], (*mapped_mat)[i]
				swaps++
			}
		}
	}
}

func CreateFronts(mat [][]int, mapped_mat [][]float64) [][][]int {
	index := 0
	result := [][][]int{}
	result = append(result, [][]int{})
	result[index] = append(result[index], mat[0])

	for i := 1; i < len(mat); i++ {
		if EuclideanLTe(mapped_mat[i], mapped_mat[i-1]) {
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

func CrossoverAndMutation(selected, population [][]int, p_crossover, p_mutation float64) [][]int {
	result := [][]int{}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	population_num := float64(len(population))

	for crossovers := 0.0; crossovers < population_num*p_crossover; {
		random_p1 := r1.Intn(len(selected))
		// random_p2 := r1.Intn(len(selected))
		//
		// for random_p2 == random_p1 {
		// 	random_p2 = r1.Intn(len(selected))
		// }
		//
		// children := Crossover(selected[random_p1], selected[random_p2])
		//
		// for _, item := range children {
		// 	if SliceValid(item, population) {
		// 		population = append(population, item)
		// 		result = append(result, item)
		// 		crossovers++
		// 	}
		// }

		child := Crossover2(selected[random_p1])

		if SliceValid(child, population) {
			population = append(population, child)
			result = append(result, child)
			crossovers++
		}

	}

	for mutations := 0.0; mutations < population_num*p_mutation; {
		random_p := r1.Intn(len(selected))

		child := MutateOnePoint(selected[random_p])
		child = MutateOnePoint(child)

		if SliceValid(child, population) {
			population = append(population, child)
			result = append(result, child)
			mutations++
		}

	}

	return result
}

func NSGA2(problem Problem) []int {
	// fmt.Println("checkpoint 1")
	obj_func := []func([]int) float64{problem.ObjectFunction1, problem.ObjectFunction2}
	population := InitializePopulation(problem.population_size, problem.problem_size)
	mapped_mat := MapMatWithObjectiveFunc(population, obj_func)
	SliceFastNondominatedSort(&population, &mapped_mat)
	selected := SelectParentsByRank(population, problem.population_size, problem.init_p_selected)
	children := CrossoverAndMutation(selected, population, problem.init_p_crossover, problem.init_p_mutation)
	// fmt.Println("checkpoint 2")
	best_gene := VecCopyInt(population[0])
	best_vec := VecCopyFloat(mapped_mat[0])

	for generation := 0; !problem.stopCondition(generation, selected, mapped_mat); generation++ {
		if EuclideanDistFloat(mapped_mat[0], []float64{0.0, 0.0}) < EuclideanDistFloat(best_vec, []float64{0.0, 0.0}) {
			best_gene = VecCopyInt(population[0])
			best_vec = VecCopyFloat(mapped_mat[0])
		}

		union := Union(population, children)
		mapped_mat = MapMatWithObjectiveFunc(union, obj_func)
		SliceFastNondominatedSort(&union, &mapped_mat)
		fronts := CreateFronts(union, mapped_mat)
		// fmt.Println("checkpoint 3")
		parents := [][]int{}
		front_l := -1
		for index, front := range fronts {
			if len(parents)+len(front) > problem.population_size {
				front_l = index
				break
			} else {
				parents = append(parents, front...)
			}
		}
		// fmt.Println("checkpoint 4")

		amount_left := problem.population_size - len(parents)
		if amount_left != 0 {
			parents = append(parents, fronts[front_l][:amount_left]...)
		}
		// fmt.Println("checkpoint 5")

		selected = SelectParentsByRank(parents, problem.population_size, problem.init_p_selected)
		population = parents
		children = CrossoverAndMutation(selected, population, problem.init_p_crossover, problem.init_p_mutation)
	}
	// fmt.Println("checkpoint 6")

	fmt.Println("best vector objective length: ", EuclideanDistFloat(best_vec, []float64{0.0, 0.0}))
	return best_gene
}
