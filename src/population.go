package src

import (
	"math/rand"
	"sync"
)

type Population struct {
	individuals       []Individual
	totalFitnessScore float64
	mutationRate      float64
	model             Model
	popSize           int
}

func (population *Population) evolve() {
	newIndividuals := make([]Individual, population.popSize)

	for i := 0; i < len(population.individuals); i++ {
		parent1 := population.model.SelectParent(population)
		parent2 := population.model.SelectParent(population)
		offSpring, _ := population.model.Crossover(parent1, parent2)

		mutationChance := rand.Float64()
		if mutationChance <= population.mutationRate {
			offSpring, _ = population.model.Mutate(offSpring)
		}

		newIndividuals[i] = offSpring
	}

	population.totalFitnessScore = 0.0
	population.individuals = newIndividuals
}

func (population *Population) evolveParallel() {
	newIndividuals := make([]Individual, population.popSize)

	var wg sync.WaitGroup
	wg.Add(len(population.individuals))
	for i := 0; i < len(population.individuals); i++ {
		go func(index int) {
			defer wg.Done() // decrement the counter when Goroutine is done

			parent1 := population.model.SelectParent(population)
			parent2 := population.model.SelectParent(population)
			offSpring, _ := population.model.Crossover(parent1, parent2)

			mutationChance := rand.Float64()
			if mutationChance <= population.mutationRate {
				offSpring, _ = population.model.Mutate(offSpring)
			}

			newIndividuals[index] = offSpring
		}(i) // passing i as an argument to the Goroutine
	}
	wg.Wait() // block until all Goroutines finish

	population.totalFitnessScore = 0.0
	population.individuals = newIndividuals
}

func (population *Population) calculateBestIndividual() Individual {
	bestFitnessScore := 0.0
	var bestIndividual Individual
	for _, individual := range population.individuals {
		currentFitnessScore := individual.CalculateFitness()
		if currentFitnessScore > bestFitnessScore {
			bestIndividual = individual
			bestFitnessScore = currentFitnessScore
		}
	}

	return bestIndividual
}

func (population *Population) getTotalFitnessScore() float64 {
	totalFitnessScore := population.totalFitnessScore
	if totalFitnessScore != 0 {
		return totalFitnessScore
	}

	for _, individual := range population.individuals {
		totalFitnessScore += individual.CalculateFitness()
	}
	population.totalFitnessScore = totalFitnessScore

	return totalFitnessScore
}
