package main

import (
	"strings"
)

type Executor struct {
	env    *TodoTXT
	source chan *RecurTask
}

func (me *Executor) ExecuteAll(target chan *RecurTask) {
	for task := range me.source {
		exec, ok := executors[strings.ToLower(task.Type)]

		if ok {
			logger.Printf("Ready to execute [task=%v]", task)
			update, err := exec(me.env, task)

			if err == nil {
				if update {
					logger.Printf("Task has been updated [%v]", task)
					target <- task
				}
			} else {
				logger.Fatal(err)
				continue
			}
		} else {
			logger.Fatalf("Recurring task [type=%s] is not supported", task.Type)
		}
	}

	close(target)
}

func NewExecutor(env *TodoTXT, source chan *RecurTask) *Executor {
	return &Executor{
		env:    env,
		source: source,
	}
}

type executor func(*TodoTXT, *RecurTask) (bool, error)

var executors = make(map[string]executor)

func init() {
	executors["weekly"] = weekly
}
