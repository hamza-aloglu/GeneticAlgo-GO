package src

import (
	"errors"
	"math/rand"
	"reflect"
)

type PermutationModel struct {
	DefaultModel
}

// Crossover is ordered crossover implementation
func (pm PermutationModel) Crossover(parent1 Individual, parent2 Individual) (Individual, error) {
	p1 := reflect.ValueOf(parent1)
	p2 := reflect.ValueOf(parent2)
	n := p1.Len()
	child := reflect.MakeSlice(p1.Type(), n, n)

	if p1.Kind() != reflect.Slice || p2.Kind() != reflect.Slice {
		return nil, errors.New("parent(s) are not slice")
	}
	if p1.Len() != p2.Len() {
		return nil, errors.New("both slices must have the same length")
	}

	bound1, bound2 := rand.Intn(n), rand.Intn(n)
	if bound1 > bound2 {
		bound1, bound2 = bound2, bound1
	}

	for i := bound1; i < bound2; i++ {
		child.Index(i).Set(p1.Index(i))
	}

	parent2Index := bound2
	childIndex := bound2
	for count := 0; count < n; count++ {
		val := p2.Index(parent2Index % n).Interface()
		if !contains(child, val) {
			child.Index(childIndex % n).Set(p2.Index(parent2Index % n))
			childIndex++
		}
		parent2Index++
	}

	return child.Interface().(Individual), nil
}

func contains(slice reflect.Value, item interface{}) bool {
	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// Mutate is swap mutation implementation
func (pm PermutationModel) Mutate(individual Individual) (Individual, error) {
	individualReflection := reflect.ValueOf(individual)

	// Pick two random 2 genes and swap their positions
	pos1, pos2 := rand.Intn(individualReflection.Len()), rand.Intn(individualReflection.Len())
	tmp := individualReflection.Index(pos1).Interface()
	individualReflection.Index(pos1).Set(individualReflection.Index(pos2))
	individualReflection.Index(pos2).Set(reflect.ValueOf(tmp))

	return individualReflection.Interface().(Individual), nil
}
