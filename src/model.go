package src

import (
	"errors"
	"math/rand"
	"reflect"
)

type Model interface {
	SelectParent(population *Population) Individual
	Crossover(parent1 Individual, parent2 Individual) (Individual, error)
	Mutate(individual Individual) (Individual, error)
}

type DefaultModel struct {
}

func (dm DefaultModel) SelectParent(population *Population) Individual {
	individuals := population.individuals
	totalFitnessScore := population.getTotalFitnessScore()
	fitnessThreshold := rand.Intn(int(totalFitnessScore))
	currentFitness := 0.0
	for _, individual := range individuals {
		currentFitness += individual.CalculateFitness()
		if currentFitness >= float64(fitnessThreshold) {
			return individual
		}
	}

	return individuals[0]
}

// Crossover is fixed point crossover. It does not ensure uniqueness of genes.
func (dm DefaultModel) Crossover(parent1 Individual, parent2 Individual) (Individual, error) {
	p1 := reflect.ValueOf(parent1)
	p2 := reflect.ValueOf(parent2)
	child := reflect.MakeSlice(p1.Type(), p1.Len(), p1.Len())

	if p1.Kind() != reflect.Slice || p2.Kind() != reflect.Slice {
		return nil, errors.New("parent(s) are not slice")
	}
	if p1.Len() != p2.Len() {
		return nil, errors.New("both slices must have the same length")
	}

	crossoverPoint := rand.Intn(p1.Len())
	for i := 0; i < p1.Len(); i++ {
		if i < crossoverPoint {
			child.Index(i).Set(p1.Index(i))
		} else {
			child.Index(i).Set(p2.Index(i))
		}
	}

	return child.Interface().(Individual), nil
}

func (dm DefaultModel) Mutate(individual Individual) (Individual, error) {
	/*allelesReflection := reflect.ValueOf(individual.GetAlleles())
	if allelesReflection.Kind() != reflect.Slice {
		return nil, errors.New("Allele must be type of slice")
	}
	individualReflection := reflect.ValueOf(individual)

	individualReflection.
		Index(rand.Intn(individualReflection.Len())).
		Set(allelesReflection.Index(rand.Intn(allelesReflection.Len())))

	return individualReflection.Interface().(Individual), nil*/
	panic("Mutate is not implemented in DefaultModel!")
}
