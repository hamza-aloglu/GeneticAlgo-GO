package main

import (
	"errors"
	"github.com/hamza-aloglu/GeneticAlgo-Go/src"
	"math"
	"math/rand"
	"time"
)

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

type Schedule struct {
	Genes                   []Day     `json:"days"`
	ScheduleStartDate       time.Time `json:"scheduleStartDate"`
	ScheduleEndDate         time.Time `json:"scheduleEndDate"`
	Tasks                   []Task    `json:"tasks"`
	IsDeadlineMustMeet      bool      `json:"isDeadlineMustMeet"`
	IndividualInsertChance  float64   `json:"individualInsertChance"`
	DayAmount               int       `json:"dayAmount"`
	AverageDifficultyPerDay float64   `json:"averageDifficultyPerDay"`

	WeightPriority   float64 `json:"weightPriority"`
	WeightDifficulty float64 `json:"weightDifficulty"`
	WeightDeadline   float64 `json:"weightDeadline"`
	WeightParentTask float64 `json:"weightParentTask"`
}

const MINIMUM_FITNESS = 0.01

func (s Schedule) CalculateFitness() float64 {
	var fitnessScore float64
	for dayIndex, day := range s.Genes {
		dayDate := s.ScheduleStartDate.AddDate(0, 0, dayIndex)
		difficultyForDay := 0
		dayTasksAmount := float64(len(day.Tasks))
		violateParentTaskAmount := 0.0
		distanceToEnd := float64(s.DayAmount-dayIndex) / float64(s.DayAmount)
		for _, task := range day.Tasks {
			difficultyForDay += task.Difficulty

			// place high priority task early in schedule
			fitnessScore += s.WeightPriority * float64(task.Priority) * distanceToEnd

			// punish not meeting deadline
			millisecondsExceedDeadline := dayDate.Sub(task.Deadline)
			daysExceedDeadline := int(millisecondsExceedDeadline.Hours() / 24)
			if daysExceedDeadline > 0 {
				daysExceedDeadline = int(math.Min(float64(daysExceedDeadline), 10))
				fitnessScore -= s.WeightDeadline * float64(daysExceedDeadline)
			} else if s.IsDeadlineMustMeet {
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
		fitnessScore -= s.WeightParentTask * violateParentTaskAmount

		// punish exceeding average difficulty
		if dayTasksAmount > 0 {
			dayAvgDifficulty := float64(difficultyForDay) / dayTasksAmount
			difficultyImpact := gaussianReward(float64(difficultyForDay), s.AverageDifficultyPerDay)
			fitnessScore += s.WeightDifficulty * difficultyImpact * dayTasksAmount * dayAvgDifficulty
		}

	}

	if fitnessScore <= 0 {
		fitnessScore = MINIMUM_FITNESS
	}

	return fitnessScore
}

func (s Schedule) GenerateIndividual() src.Individual {
	var localTasks []Task
	localTasks = append(localTasks, s.Tasks...)
	individual := generateSchedule()
	for len(localTasks) > 0 {
		for i := 0; i < s.DayAmount && len(localTasks) > 0; i++ {
			if rand.Float64() < s.IndividualInsertChance {
				// insert randomly selected task to the day
				randomTaskIndex := rand.Intn(len(localTasks))
				individual.Genes[i].Tasks = append(individual.Genes[i].Tasks, localTasks[randomTaskIndex])

				// remove randomly selected task from whole tasks
				localTasks = append(localTasks[:randomTaskIndex], localTasks[randomTaskIndex+1:]...)
			}
		}
	}
	return individual
}

func main() {
	defer timer("main")()

	ga := src.NewCustomGA(20, 1000, 0.3, generateSchedule(), ScheduleModel{})
	var bestSchedule Schedule
	bestSchedule = ga.RunWithLog(printSchedule).(Schedule)

	println(bestSchedule.Genes)

	printBenchmark()
}

func (s Schedule) findDayOfTaskByTitle(title string) (int, error) {
	for dayNumber, day := range s.Genes {
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
