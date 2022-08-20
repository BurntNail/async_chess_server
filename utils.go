package main

import "math"

type JVec[T comparable] struct {
	data []T
}

func (jv *JVec[T]) Add(item T) {
	jv.data = append(jv.data, item)
}

// Returns true if the element was removed, false if it wasn't
func (jv *JVec[T]) Remove(item T) bool {
	if ind := jv.IndexOf(item); ind != -1 {
		jv.data[ind] = jv.data[len(jv.data)-1]
		jv.data = jv.data[:len(jv.data)-1]
		return true
	} else {
		return false
	}
}

// Returns -1 if not found
func (jv *JVec[T]) IndexOf(item T) int {
	ind := -1
	for index, el := range jv.data {
		if el == item {
			ind = index
			break
		}
	}

	return ind
}

func AbsInt(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func PythagDist(dx, dy int) float64 {
	x := float64(dx * dx)
	y := float64(dy * dy)
	return math.Sqrt(x + y)
}
