package nsga2

import (
	"math"
	"math/rand"
	"time"
)

func SliceEqual(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func SliceIndexInMat(slice []int, mat [][]int) int {
	for index, item := range mat {
		if SliceEqual(slice, item) {
			return index
		}
	}
	return -1
}

func EuclideanDist(a, b []int) float64 {
	return math.Sqrt((math.Pow(float64(a[0]), 2.0) - math.Pow(float64(b[0]), 2.0)) + (math.Pow(float64(a[1]), 2.0) - math.Pow(float64(b[1]), 2.0)))
}

func FuncsOnSlice(slice []int, obj_funcs []func([]int) float64) []float64 {
	output := []float64{}

	for _, obj_func := range obj_funcs {
		output = append(output, obj_func(slice))
	}

	return output
}

func SliceLTE(a, b []float64) bool {
	sum := 0
	for i := 0; i < len(a); i++ {
		con := b[i] - a[i]
		switch {
		case con > 0.0:
			sum++
		case con < 0.0:
			sum--
		}
	}
	return sum >= 0
}

func Crossover(item []int) []int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rand_index := r1.Intn(len(item) - 1)

	child1 := make([]int, len(item)-rand_index)
	child2 := make([]int, rand_index)
	copy(child1, item[rand_index:])
	copy(child2, item[:rand_index])

	return append(child1, child2...)
}

func Mutate(item []int) []int {
	child := make([]int, len(item))
	copy(child, item)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rand_index1 := r1.Intn(len(item))
	rand_index2 := r1.Intn(len(item))

	child[rand_index1], child[rand_index2] = child[rand_index2], child[rand_index1]

	return child
}

// func Crossover1(item1, item2 []int) []int, []int {
//   child1, child2 := make([]int, len(item1)), make([]int, len(item2))
//
// }
