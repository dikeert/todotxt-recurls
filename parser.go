package main

import (
	"errors"
	"fmt"
)

type RecurParser struct {
	source chan *RecurRecord
}

func (me *RecurParser) ParseAll(target chan *RecurTask) {
	for record := range me.source {
		logger.Printf("Received [line=%s] to parse, parsing", record.Line)
		if task, err := parse(record.Line); err == nil {
			task.ID = record.ID
			logger.Printf("Parsed [task=%v]", task)
			target <- task
		} else {
			logger.Fatal(err)
		}
	}

	close(target)
}

func (me *RecurParser) MarshallAll(source chan *RecurTask, target chan *RecurRecord) {
	for task := range source {
		target <- marshall(task)
	}

	close(target)
}

func NewParser(source chan *RecurRecord) *RecurParser {
	return &RecurParser{
		source: source,
	}
}

type state int

const (
	cmd state = iota
	arg
	attr
	value
	task
	final
)

type token int

const (
	space token = iota
	colon
	comma
	equals
	any
	none
)

type listener func(state, state, string, *RecurTask) error

type rule struct {
	current  state
	symbol   token
	next     state
	listener listener
}

var rules = []*rule{
	{cmd, comma, attr, onCmd},
	{cmd, colon, task, onCmd},
	{cmd, space, arg, onCmd},
	{attr, equals, value, attrs.onAttr},
	{value, comma, attr, attrs.onValue},
	{value, space, arg, attrs.onValue},
	{value, colon, task, attrs.onValue},
	{arg, space, arg, onArg},
	{arg, colon, task, onArg},
	{task, none, final, onTask},
}

func marshall(task *RecurTask) *RecurRecord {
	attrs := ""
	for k, v := range task.Attr {
		if attrs != "" {
			attrs += ","
		}

		attrs += fmt.Sprintf("%s=%s", k, v)
	}

	args := ""
	for _, v := range task.Args {
		if args != "" {
			args += " "
		}
		args += v
	}

	line := task.Type
	if attrs != "" {
		line += fmt.Sprintf(",%s", attrs)
	}
	if args != "" {
		line += fmt.Sprintf(" %s", args)
	}

	return &RecurRecord{
		ID:   task.ID,
		Line: fmt.Sprintf("%s:%s", line, task.Todo),
	}
}

func parse(line string) (*RecurTask, error) {
	task := NewTask()
	state := cmd

	buff := ""
	for _, r := range line {
		ch := string(r)

		if next, reset, err := transition(buff, ch, state, task); err == nil {
			if reset == true {
				state = next
				buff = ""
			} else {
				buff += ch
			}
		} else {
			return nil, err
		}
	}

	if done, err := exit(buff, state, task); err == nil {
		if done {
			return task, nil
		} else {
			return nil, errors.New("Fail to parse Recurring Task line")
		}
	} else {
		return nil, err
	}
}

func transition(buff string, ch string, state state, task *RecurTask) (state, bool, error) {
	for _, rule := range rules {
		if rule.current == state && rule.symbol == sym(ch) {
			if rule.listener != nil {
				if err := rule.listener(rule.current, rule.next, buff, task); err != nil {
					return state, false, err
				}
			}

			return rule.next, true, nil
		}
	}

	return state, false, nil
}

func exit(buff string, state state, task *RecurTask) (bool, error) {
	for _, rule := range rules {
		if rule.current == state && rule.next == final {
			if rule.listener != nil {
				if err := rule.listener(rule.current, rule.next, buff, task); err != nil {
					return false, err
				}
			}

			return true, nil
		}
	}

	return false, nil
}

func sym(ch string) token {
	switch ch {
	case " ":
		return space
	case ",":
		return comma
	case "=":
		return equals
	case ":":
		return colon
	}

	return any
}

func onCmd(curr state, next state, cmd string, task *RecurTask) error {
	if len(cmd) < 1 {
		return errors.New("Recurring task type can't be empty")
	}

	task.Type = cmd
	return nil
}

func onArg(curr state, next state, arg string, task *RecurTask) error {
	if len(arg) < 1 {
		return errors.New("Argument can't be empty")
	}

	task.Args = append(task.Args, arg)
	return nil
}

func onTask(curr state, next state, todo string, task *RecurTask) error {
	if len(todo) < 1 {
		return errors.New("Todo can't be empty for recurring task")
	} else {
		task.Todo = todo
		return nil
	}
}

type attrListener struct {
	attr string
}

func (me *attrListener) onAttr(curr state, next state, attr string, task *RecurTask) error {
	me.attr = attr
	return nil
}

func (me *attrListener) onValue(curr state, next state, value string, task *RecurTask) error {
	task.Attr[me.attr] = value
	return nil
}

var attrs = &attrListener{}
