package main

import (
	"errors"
	"fmt"
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	alleles []string = strings.Split("abcdefghijklmnopqrstuvwxyz,! ", "")
	target  []string = strings.Split("hello, world!", "")
)

// Characters represent genes and implements Individual
type Characters []string

func (c Characters) CalculateFitness() float64 {
	result := 0.0
	for i, gene := range c {
		if gene == target[i] {
			result += 1.0
		}
	}

	return result
}

func (c Characters) GenerateIndividual() src.Individual {
	individual := make(Characters, len(target))
	for i := 0; i < len(target); i++ {
		individual[i] = alleles[rand.Intn(len(alleles))]
	}
	return individual
}

type StringModel struct {
	src.DefaultModel
}

func (sm StringModel) Mutate(individual src.Individual) (src.Individual, error) {
	// Assert the individual to type Characters
	charactersIndividual, ok := individual.(Characters)
	if !ok {
		// Handle the case where the assertion fails, e.g., return an error
		return nil, errors.New("individual is not of type Characters")
	}
	// Perform the mutation on charactersIndividual
	charactersIndividual[rand.Intn(len(charactersIndividual))] = alleles[rand.Intn(len(alleles))]
	return charactersIndividual, nil
}

func main() {
	defer timer("main")()

	// Models can be given using Enum. (Models.StringModel).
	ga := src.NewCustomGA(40, 2500, 0.5, Characters{}, StringModel{})
	var bestCharacters Characters
	bestCharacters = ga.Run().(Characters)
	for _, gen := range bestCharacters {
		fmt.Printf("\n%v", gen)
	}

	// Collect and print memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory: %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("Total allocated memory: %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System memory obtained from OS: %v MB\n", m.Sys/1024/1024)

	// Collect and print CPU usage
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getpid()), "-o", "%cpu")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error getting CPU usage:", err)
	}
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
