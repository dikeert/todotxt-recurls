package main

import (
	"os"
)

func main() {
	env := TodoEnviron()

	if r, err := OpenRecur(env); err == nil {
		lines := make(chan *RecurRecord)
		tasks := make(chan *RecurTask)
		updates := make(chan *RecurTask)
		writes := make(chan *RecurRecord)

		p := NewParser(lines)
		e := NewExecutor(env, tasks)

		go p.ParseAll(tasks)
		go e.ExecuteAll(updates)
		go p.MarshallAll(updates, writes)
		go r.WriteAll(writes)

		r.ReadAll(lines)
	} else {
		logger.Fatal(err)
	}

	if cmd, err := CreateCmd(); err != nil {
		logger.Fatal(os.Stderr, err)
		os.Exit(1)
	} else if err := cmd.Run(); err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
