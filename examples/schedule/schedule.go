package main

import (
	"errors"
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"math"
	"math/rand"
	"time"
)

// api input
var scheduleStartDate = time.Date(2024, time.March, 5, 0, 0, 0, 0, time.UTC)
var scheduleEndDate = time.Date(2024, time.March, 12, 0, 0, 0, 0, time.UTC)
var csvRecords = readCsvFile("/Users/hamza/Projects/Go/GeneticAlgo-Go/examples/schedule/tasks.csv")
var slicedCsvRecords = csvRecords[1:]
var tasks = convertIntoTasks(slicedCsvRecords)
var isDeadlineMustMeet = false

type Task struct {
	Title           string    `json:"title"`
	Deadline        time.Time `json:"deadline"`
	Difficulty      int       `json:"difficulty"`
	Priority        int       `json:"priority"`
	ParentTaskTitle string    `json:"parentTaskTitle"`
}

type Day struct {
	Tasks []Task `json:"tasks"`
}

type Schedule []Day

// internal configs
var duration = scheduleEndDate.Sub(scheduleStartDate)
var dayAmount = int(duration.Hours() / 24)
var insertChance = 0.2
var totalDifficulty = sumDifficulty(tasks)
var averageDifficultyPerDay = totalDifficulty / dayAmount

const MINIMUM_FITNESS = 0.01

var weightDeadline = 10.0
var weightDifficulty = 2.0
var weightPriority = 2.0
var weightParentTask = 10.0

func (s Schedule) CalculateFitness() float64 {
	var fitnessScore float64
	for dayIndex, day := range s {
		dayDate := scheduleStartDate.AddDate(0, 0, dayIndex)
		difficultyForDay := 0
		dayTasksAmount := float64(len(day.Tasks))
		violateParentTaskAmount := 0.0
		distanceToEnd := float64(dayAmount-dayIndex) / float64(dayAmount)
		for _, task := range day.Tasks {
			difficultyForDay += task.Difficulty

			// place high priority task early in schedule
			fitnessScore += weightPriority * float64(task.Priority) * distanceToEnd

			// punish not meeting deadline
			millisecondsExceedDeadline := dayDate.Sub(task.Deadline)
			daysExceedDeadline := int(millisecondsExceedDeadline.Hours() / 24)
			if daysExceedDeadline > 0 {
				daysExceedDeadline = int(math.Min(float64(daysExceedDeadline), 10))
				fitnessScore -= weightDeadline * float64(daysExceedDeadline)
			} else if isDeadlineMustMeet && daysExceedDeadline <= 0 {
				return MINIMUM_FITNESS
			}

			// calculate parent task violation
			if task.ParentTaskTitle != "" {
				parentDayIndex, err := s.findDayOfTaskByTitle(task.ParentTaskTitle)
				if err != nil || parentDayIndex > dayIndex {
					violateParentTaskAmount++
				}
			}
		}

		// punish if parent task is not scheduled before current task
		violateParentTaskAmount = math.Min(float64(violateParentTaskAmount), 10)
		fitnessScore -= weightParentTask * violateParentTaskAmount

		// punish exceeding average difficulty
		if dayTasksAmount > 0 {
			dayAvgDifficulty := float64(difficultyForDay) / dayTasksAmount
			difficultyImpact := gaussianReward(float64(difficultyForDay), float64(averageDifficultyPerDay))
			fitnessScore += weightDifficulty * difficultyImpact * dayTasksAmount * dayAvgDifficulty
		}

	}

	if fitnessScore <= 0 {
		fitnessScore = MINIMUM_FITNESS
	}

	return fitnessScore
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
				individual[i].Tasks = append(individual[i].Tasks, localTasks[randomTaskIndex])

				// remove randomly selected task from whole tasks
				localTasks = append(localTasks[:randomTaskIndex], localTasks[randomTaskIndex+1:]...)
			}
		}
	}
	return individual
}

func main() {
	defer timer("main")()

	ga := src.NewCustomGA(100, 1000, 0.3, Schedule{}, ScheduleModel{})
	var bestSchedule Schedule
	bestSchedule = ga.RunWithLog(printSchedule).(Schedule)

	println(bestSchedule)

	printBenchmark()
}

func (s Schedule) findDayOfTaskByTitle(title string) (int, error) {
	for dayNumber, day := range s {
		if day.containsTaskByTitle(title) {
			return dayNumber, nil
		}
	}

	return -1, errors.New("no task with title: " + title)
}

func (d Day) containsTaskByTitle(title string) bool {
	for _, task := range d.Tasks {
		if task.Title == title {
			return true
		}
	}
	return false
}

// reward based on Gaussian function
func gaussianReward(difficultyForDay, averageDifficultyPerDay float64) float64 {
	distance := difficultyForDay - averageDifficultyPerDay
	sigma := 1.7 /* set your desired sigma value */
	exponent := -0.022 * (distance / sigma) * (distance / sigma) * (distance / sigma)
	return math.Exp(exponent)
}
