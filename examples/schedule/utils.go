package main

import (
	"encoding/csv"
	"fmt"
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

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

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func convertIntoTasks(input [][]string) []Task {
	tasks := []Task{}
	for _, row := range input {
		newTask := Task{}
		newTask.Deadline, _ = time.Parse(time.DateOnly, row[0])
		newTask.Difficulty, _ = strconv.Atoi(row[1])
		newTask.Priority, _ = strconv.Atoi(row[2])
		newTask.Title = row[3]
		tasks = append(tasks, newTask)
	}

	return tasks
}

func sumDifficulty(tasks []Task) int {
	totalDifficulty := 0
	for _, task := range tasks {
		totalDifficulty += task.Difficulty
	}
	return totalDifficulty
}

func printSchedule(individual src.Individual) {
	schedule := individual.(Schedule)
	for dayCount, day := range schedule {
		fmt.Printf("DAY - %v\n", dayCount)
		for _, task := range day.Tasks {
			println(task.Title)
			println(task.Deadline.String())
			println(task.Priority)
			println(task.Difficulty)
		}
		println("------------------------------------------------------------------")
	}
	println()
	println()
}
