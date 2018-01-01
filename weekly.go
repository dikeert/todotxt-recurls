package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func weekly(env *TodoTXT, task *RecurTask) (bool, error) {
	logger.Printf("Received weekly recurring task [args=%v]", task.Args)

	today := today()
	weekday := weekday(today)

	if len(task.Args) < 1 {
		return false, errors.New(fmt.Sprint("Not enough arguments for task [type=%s]", task.Type))
	}

	if skip, err := skip(today, task); err == nil {
		if skip {
			return false, nil
		}
	} else {
		return false, err
	}

	for _, day := range task.Args {
		if strings.ToLower(day) == strings.ToLower(weekday) {
			logger.Printf("Current [weekday=%s] matches task [weekday=%s], executing...", weekday, day)
			env.AddTodo(task.Todo)

			update(&today, task)
			return true, nil
		}
	}

	return false, nil
}

const LAST = "last"
const LAYOUT = "02-01-2006"

func today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func update(today *time.Time, task *RecurTask) {
	update := today.Format(LAYOUT)
	logger.Printf("Calculated [value=%s] for [attr=%s]", update, LAST)
	task.Attr[LAST] = update
}

func weekday(today time.Time) string {
	var weekday string

	switch today.Weekday() {
	case time.Monday:
		weekday = "Monday"
	case time.Tuesday:
		weekday = "Tuesday"
	case time.Wednesday:
		weekday = "Wednesday"
	case time.Thursday:
		weekday = "Thursday"
	case time.Friday:
		weekday = "Friday"
	case time.Saturday:
		weekday = "Saturday"
	case time.Sunday:
		weekday = "Sunday"
	}

	logger.Printf("Resolved current [weekday=%s]", weekday)

	return weekday
}

func skip(today time.Time, task *RecurTask) (bool, error) {
	if len(task.Attr) > 0 {
		if lastExec, err := time.Parse(LAYOUT, task.Attr[LAST]); err == nil {
			logger.Printf("Last execution [date=%v]", lastExec)

			if !lastExec.Before(today) {
				logger.Printf("Task [%v] has already been processed [today=%v], skipping", task, today)
				return true, nil
			} else {
				logger.Printf(
					"Task [%v] last execution [date=%v] hasn't been executed [today=%v], executing",
					task, lastExec, today)
			}
		} else {
			return true, err
		}
	}

	return false, nil
}
