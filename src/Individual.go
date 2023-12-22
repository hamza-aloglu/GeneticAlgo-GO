package src

type Individual interface {
	CalculateFitness() float64
	GenerateIndividual() Individual
}
