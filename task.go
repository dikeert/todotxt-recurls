package main

type RecurTask struct {
	ID   uint
	Type string
	Args []string
	Todo string
	Attr map[string]string
}

func NewTask() *RecurTask {
	task := &RecurTask{}
	task.Attr = make(map[string]string)
	return task
}
