package src

import (
	"fmt"
	"reflect"
	"sync"
)

type GA struct {
	generationNumber int
	populationSize   int
	population       Population
}

func (g *GA) Run() Individual {
	for i := 0; i < g.generationNumber; i++ {
		g.population.evolveParallel()
		printBestOne(g)
	}

	return g.population.calculateBestIndividual()
}

// maybe use builder pattern?
func NewDefaultGA(generationNumber int, populationSize int, mutationRate float64, individual Individual) GA {
	return GA{
		generationNumber: generationNumber,
		population: Population{
			mutationRate: mutationRate,
			individuals:  generateInitialIndividuals(individual.GenerateIndividual, populationSize),
			model:        PermutationModel{},
			popSize:      populationSize,
		},
	}
}

func NewCustomGA(generationNumber int, populationSize int, mutationRate float64, individual Individual, model Model) GA {
	// model = createModel(modelType)
	// --inside the createModel--
	// return modelFactories[modelType]()
	// eğer bir config dosyasından bir şeyler almam gerekse yukarıdaki fonksiyona parametre olarak yollayabilirim.

	return GA{
		generationNumber: generationNumber,
		population: Population{
			mutationRate: mutationRate,
			individuals:  generateInitialIndividuals(individual.GenerateIndividual, populationSize),
			model:        model,
			popSize:      populationSize,
			fitnessCache: FitnessCache{
				cache: sync.Map{},
			},
			totalFitnessScore: 0.0,
		},
	}
}

func generateInitialIndividuals(generateIndividual func() Individual, populationSize int) []Individual {
	initialIndividuals := make([]Individual, populationSize)
	for i := 0; i < populationSize; i++ {
		initialIndividuals[i] = generateIndividual()
	}

	return initialIndividuals
}

func printBestOne(g *GA) {
	bestOne := g.population.calculateBestIndividual()
	fmt.Printf("\nfitness of best one: %v \n", bestOne.CalculateFitness())

	bestOneReflection := reflect.ValueOf(bestOne)
	if bestOneReflection.Kind() == reflect.Slice {
		for i := 0; i < bestOneReflection.Len(); i++ {
			fmt.Printf(" %v ", bestOneReflection.Index(i))
		}
	}
}
