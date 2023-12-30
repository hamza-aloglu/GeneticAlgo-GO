package main

import (
	"fmt"
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Schedule -> days -> tasks

type Task struct {
	Title      string
	Deadline   time.Time
	Difficulty int
	Priority   int
}

type Day struct {
	tasks []Task
}

type Schedule []Day

var dayAmount = 5
var insertChance = 0.2

// Creating dummy data for Task
var task1 = Task{
	Title:      "Implement User Authentication",
	Deadline:   time.Now().AddDate(0, 0, 7), // Deadline is set to 7 days from now
	Difficulty: 3,
	Priority:   2,
}

var task2 = Task{
	Title:      "Refactor Database Layer",
	Deadline:   time.Now().AddDate(0, 0, 14), // Deadline is set to 14 days from now
	Difficulty: 4,
	Priority:   1,
}

var task3 = Task{
	Title:      "Update UI Components",
	Deadline:   time.Now().AddDate(0, 0, 5), // Deadline is set to 5 days from now
	Difficulty: 2,
	Priority:   3,
}

var tasks = []Task{task1, task2, task3}

func (s Schedule) CalculateFitness() float64 {
	panic("implement calcualte fitness")
}

func (s Schedule) GenerateIndividual() src.Individual {
	var localTasks []Task
	localTasks = append(localTasks, tasks...)
	individual := make(Schedule, dayAmount)
	for len(localTasks) > 0 {
		for i := 0; i < dayAmount && len(localTasks) > 0; i++ {
			if rand.Float64() < insertChance {
				// insert randomly selected task to the day
				randomTaskIndex := rand.Intn(len(localTasks))
				individual[i].tasks = append(individual[i].tasks, localTasks[randomTaskIndex])

				// remove randomly selected task from whole tasks
				localTasks = append(localTasks[:randomTaskIndex], localTasks[randomTaskIndex+1:]...)
			}
		}
	}
	return individual
}

func main() {
	defer timer("main")()

	ga := src.NewDefaultGA(40, 2500, 0.5, Schedule{})
	var bestSchedule Schedule
	bestSchedule = ga.Run().(Schedule)
	for _, gen := range bestSchedule {
		fmt.Printf("\n%v", gen)
	}

	printBenchmark()
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func printBenchmark() {
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
