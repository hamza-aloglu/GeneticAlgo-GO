package src

type GA struct {
	generationNumber int
	populationSize   int
	population       Population
}

type printIndividual func(individual Individual)

// Implement user provided callback function that takes individual as parameter
// We will call this callback with parameter g.population.calculateBestIndividual
// User can get properties of this individual and save it in an array and use it.
func (g *GA) Run() Individual {
	for i := 0; i < g.generationNumber; i++ {
		g.population.evolveParallel()
	}

	return g.population.calculateBestIndividual()
}

func (g *GA) RunWithLog(printIndividual printIndividual) Individual {
	for i := 0; i < g.generationNumber; i++ {
		g.population.evolveParallel()
		bestOne := g.population.calculateBestIndividual()
		printIndividual(bestOne)
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

func NewCustomGA(generationNumber int, populationSize int, mutationRate float64, elitismRate float64, individual Individual, model Model) GA {
	// model = createModel(modelType)
	// --inside the createModel--
	// return modelFactories[modelType]()
	// eğer bir config dosyasından bir şeyler almam gerekse yukarıdaki fonksiyona parametre olarak yollayabilirim.

	return GA{
		generationNumber: generationNumber,
		population: Population{
			mutationRate:      mutationRate,
			individuals:       generateInitialIndividuals(individual.GenerateIndividual, populationSize),
			model:             model,
			popSize:           populationSize,
			totalFitnessScore: 0.0,
			elitismRate:       elitismRate,
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
