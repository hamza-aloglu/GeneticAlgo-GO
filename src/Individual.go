package src

import (
	"fmt"
	"sync"
)

type Individual interface {
	CalculateFitness() float64
	GenerateIndividual() Individual
}

type FitnessCache struct {
	cache sync.Map
}

// It should only read, not write.
func (fc *FitnessCache) GetFitness(individual Individual) float64 {
	score, ok := fc.cache.Load(getUniqueID(individual))
	if ok {
		return score.(float64)
	}

	score = individual.CalculateFitness()
	fc.cache.Store(getUniqueID(individual), score)
	return score.(float64)
}

func (fc *FitnessCache) Clear() {
	fc.cache = sync.Map{}
}

// Helper function to get a unique ID for an individual
func getUniqueID(individual Individual) string {
	return fmt.Sprintf("%p", individual)
}
