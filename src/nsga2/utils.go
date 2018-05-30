package nsga2

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/Arafatk/glot"
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

func FuncsOnSlice(slice []int, obj_funcs []func([]int) float64) []float64 {
	output := []float64{}

	for _, obj_func := range obj_funcs {
		output = append(output, obj_func(slice))
	}

	return output
}

// func LT(a, b []float64) bool {
// 	sum := 0
// 	for i := 0; i < len(a); i++ {
// 		con := b[i] - a[i]
// 		switch {
// 		case con > 0.0:
// 			sum++
// 		case con < 0.0:
// 			sum--
// 		}
// 	}
// 	return sum > 0
// }
//
// func LTe(a, b []float64) bool {
// 	sum := 0
// 	for i := 0; i < len(a); i++ {
// 		con := b[i] - a[i]
// 		switch {
// 		case con > 0.0:
// 			sum++
// 		case con < 0.0:
// 			sum--
// 		}
// 	}
// 	return sum >= 0
// }

// func LTAbs(a, b []float64) bool {
// 	sum := 0.0
// 	for i := 0; i < len(a); i++ {
// 		sum += b[i] - a[i]
// 	}
// 	return sum > 0
// }
//
// func LTeAbs(a, b []float64) bool {
// 	sum := 0.0
// 	for i := 0; i < len(a); i++ {
// 		sum += b[i] - a[i]
// 	}
// 	return sum >= 0
// }

func EuclideanLT(a, b []float64) bool {
	return EuclideanDistFloat(a, []float64{0, 0}) < EuclideanDistFloat(b, []float64{0, 0})
}

func EuclideanLTe(a, b []float64) bool {
	return EuclideanDistFloat(a, []float64{0, 0}) <= EuclideanDistFloat(b, []float64{0, 0})
}

func Crossover(item1, item2 []int) [][]int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rand_index := r1.Intn(len(item1)-2) + 1

	child1, child2 := make([]int, rand_index), make([]int, rand_index)
	copy(child1, item1[:rand_index])
	copy(child2, item2[:rand_index])

	child1 = Append2ListsInNewInt(child1, item2[rand_index:])
	child2 = Append2ListsInNewInt(child2, item1[rand_index:])

	return [][]int{child1, child2}
}

func Crossover2(item []int) []int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rand_index := r1.Intn(len(item)-3) + 2

	child_part1 := make([]int, rand_index)
	copy(child_part1, item[:rand_index])

	child_part2 := MutateOnePoint(item[rand_index:])
	child_part2 = MutateOnePoint(item[rand_index:])
	child_part2 = MutateOnePoint(item[rand_index:])

	return Append2ListsInNewInt(child_part1, child_part2)
}

func MutateOnePoint(item []int) []int {
	child := make([]int, len(item))
	copy(child, item)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rand_index1 := r1.Intn(len(item))
	rand_index2 := r1.Intn(len(item))

	child[rand_index1], child[rand_index2] = child[rand_index2], child[rand_index1]

	return child
}

func EuclideanDist(a, b []int) float64 {
	return math.Sqrt(math.Pow(float64(a[0])-float64(b[0]), 2.0) + math.Pow(float64(a[1])-float64(b[1]), 2.0))
}

func EuclideanDistFloat(a, b []float64) float64 {
	return math.Sqrt(math.Pow(a[0]-b[0], 2.0) + math.Pow(a[1]-b[1], 2.0))
}

func AvgEuclideanDist(slice [][]float64) float64 {
	sum := 0.0
	for _, item := range slice {
		sum += EuclideanDistFloat(item, []float64{0.0, 0.0})
	}
	return sum / float64(len(slice))
}

func Append2ListsInNewInt(slice_a, slice_b []int) []int {
	new_slice := make([]int, len(slice_a)+len(slice_b))

	for i := 0; i < len(slice_a); i++ {
		new_slice[i] = slice_a[i]
	}

	for i := 0; i < len(slice_b); i++ {
		new_slice[i+len(slice_a)] = slice_b[i]
	}

	return new_slice
}

func Append2ListsInNewFloat(slice_a, slice_b []float64) []float64 {
	new_slice := make([]float64, len(slice_a)+len(slice_b))

	for i := 0; i < len(slice_a); i++ {
		new_slice[i] = slice_a[i]
	}

	for i := 0; i < len(slice_b); i++ {
		new_slice[i+len(slice_a)] = slice_b[i]
	}

	return new_slice
}

func Union(slice_a, slice_b [][]int) [][]int {
	union := [][]int{}
	for _, item := range append(slice_a, slice_b...) {
		if SliceIndexInMat(item, union) == -1 {
			item_copy := make([]int, len(item))
			copy(item_copy, item)
			union = append(union, item_copy)
		}
	}
	return union
}

func SliceValid(slice []int, mat [][]int) bool {
	slice_bool := make([]bool, len(slice))
	for _, item := range slice {
		if slice_bool[item] {
			return false
		}
		slice_bool[item] = true
	}

	// if SliceIndexInMat(slice, mat) != -1 {
	// 	return false
	// }

	return true
}

func VecCopyInt(item []int) []int {
	tmp := make([]int, len(item))
	copy(tmp, item)
	return tmp
}

func VecCopyFloat(item []float64) []float64 {
	tmp := make([]float64, len(item))
	copy(tmp, item)
	return tmp
}

func transpose(x [][]float64) [][]float64 {
	out := make([][]float64, len(x[0]))
	for i := 0; i < len(x); i += 1 {
		for j := 0; j < len(x[0]); j += 1 {
			out[j] = append(out[j], x[i][j])
		}
	}
	return out
}

func Plot(generation int, children [][]int, mapped_mat [][]float64) {
	dimensions := 2
	// The dimensions supported by the plot
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)
	pointGroupName := "Objective Vector"
	style := "points"
	points := transpose(mapped_mat)
	// Adding a point group
	plot.AddPointGroup(pointGroupName, style, points)
	// A plot type used to make points/ curves and customize and save them as an image.
	plot.SetTitle("Evolution of the Objective functions")
	// Optional: Setting the title of the plot
	plot.SetXLabel("Objective Funtion 0")
	plot.SetYLabel("Objective Funtion 1")
	// Optional: Setting label for X and Y axis
	// min_x := points[0][MinIndexInSlice(points[0])]
	// max_x := points[0][MaxIndexInSlice(points[0])]
	// min_y := points[1][MinIndexInSlice(points[1])]
	// max_y := points[1][MaxIndexInSlice(points[1])]
	//
	// margin_x := (max_x - min_x) * 0.1
	// margin_y := (max_y - min_y) * 0.1
	//
	// x_lower := int(min_x - margin_x)
	// x_upper := int(max_x - margin_x)
	// y_lower := int(min_y - margin_y)
	// y_upper := int(max_y - margin_y)

	plot.SetXrange(1000, 5000)
	plot.SetYrange(1500, 6000)
	// Optional: Setting axis ranges
	plot.SavePlot("./results/generation " + strconv.Itoa(generation) + ".png")
}

func MinIndexInSlice(slice []float64) int {
	min_index := 0

	for index, item := range slice {
		if slice[min_index] > item {
			min_index = index
		}
	}

	return min_index
}

func MaxIndexInSlice(slice []float64) int {
	max_index := 0

	for index, item := range slice {
		if slice[max_index] < item {
			max_index = index
		}
	}

	return max_index
}
