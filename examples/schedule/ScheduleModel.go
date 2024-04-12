package main

import (
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"math/rand"
)

type ScheduleModel struct {
	src.PermutationModel
}

func (sm ScheduleModel) Crossover(parent1 src.Individual, parent2 src.Individual) (src.Individual, error) {
	// Map of tasks for each parent: task -> day number. Sorted by task.
	// When selecting tasks for child, look at each task and randomly choose one and insert it into child' relevant day.
	parent1Schedule := parent1.(Schedule)
	parent2Schedule := parent2.(Schedule)
	tasksByDayParent1 := createTasksByDay(parent1Schedule)
	tasksByDayParent2 := createTasksByDay(parent2Schedule)

	child := generateSchedule()
	for _, task := range parent1Schedule.Tasks {
		if rand.Float64() <= 0.5 {
			child.Genes[tasksByDayParent1[task.Title]].Tasks = append(child.Genes[tasksByDayParent1[task.Title]].Tasks, task)
		} else {
			child.Genes[tasksByDayParent2[task.Title]].Tasks = append(child.Genes[tasksByDayParent2[task.Title]].Tasks, task)
		}
	}

	return child, nil
}

func createTasksByDay(schedule Schedule) map[string]int {
	tasksByDay := make(map[string]int, len(schedule.Tasks))
	for _, task := range schedule.Tasks {
		dayOftask, err := schedule.findDayOfTaskByTitle(task.Title)
		if err != nil {
			panic(err)
		}
		tasksByDay[task.Title] = dayOftask
	}

	return tasksByDay
}
