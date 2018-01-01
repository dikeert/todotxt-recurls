package main

import (
	"fmt"
	"testing"
)

func TestMarhalling(t *testing.T) {
	task := &RecurTask{
		ID:   0,
		Type: "Weekly",
		Todo: "do important thing @home",
		Args: []string{"Sunday", "Monday"},
		Attr: map[string]string{"last": "31-12-2017"},
	}

	record := marshall(task)

	if record.ID != 0 {
		t.Fatal("Expected to have Recurrent Record with [id=%v] but got [id=%v]", task.ID, record.ID)
		t.FailNow()
	}

	expected := "Weekly,last=31-12-2017 Sunday Monday:do important thing @home"
	if record.Line != expected {
		t.Fatal("Expected to have recurrent task [txt=%s] but got [txt=%s]", expected, record.Line)
	}
}

func TestBasicParsing(t *testing.T) {
	line := "blah foo bar:do a thing"

	if _, err := parse(line); err != nil {
		t.Errorf("Unable to parse a simple recurring task line")
		t.FailNow()
	}
}

func TestBasicNoArgParsing(t *testing.T) {
	line := "blah:do a thing"

	if _, err := parse(line); err != nil {
		t.Errorf("Unable to parse a simple recurring task line")
		t.FailNow()
	}
}

func TestBasicAttrOnlyParsing(t *testing.T) {
	line := "blah,foo=var,key=value:do a thing"

	if _, err := parse(line); err != nil {
		t.Errorf("Unable to parse a simple recurring task line")
		t.FailNow()
	}
}

func TestBasiAttrParsing(t *testing.T) {
	line := "blah,foo=bar,key=value foo bar:do a thing"

	if _, err := parse(line); err != nil {
		t.Errorf("Unable to parse a simple recurring task line")
		t.FailNow()
	}
}

func TestIsAbleToParseRecurringTaskLine(t *testing.T) {
	cases := []*aLine{
		&aLine{
			"weekly",
			"",
			"monday",
			"wednesday",
			"friday",
			"have a glass of beer @work +failed_project",
			nil,
		},

		&aLine{
			"weekly",
			"foo=bar,bazz=blah",
			"monday",
			"wednesday",
			"friday",
			"have a glass of beer @work +failed_project",
			map[string]string{
				"foo":  "bar",
				"bazz": "blah",
			},
		},
	}

	for _, val := range cases {
		doTestALine(val, t)
	}
}

type aLine struct {
	cmd       string
	attrs     string
	monday    string
	wednesday string
	friday    string
	todo      string
	result    map[string]string
}

func doTestALine(line *aLine, t *testing.T) {
	var task string

	if len(line.attrs) > 0 {
		task = fmt.Sprintf(
			"%s,%s %s %s %s:%s",
			line.cmd,
			line.attrs,
			line.monday,
			line.wednesday,
			line.friday,
			line.todo)
	} else {
		task = fmt.Sprintf(
			"%s %s %s %s:%s",
			line.cmd,
			line.monday,
			line.wednesday,
			line.friday,
			line.todo)
	}

	result, err := parse(task)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fail := false
	if result.Type != line.cmd {
		t.Errorf("Failed to parse the Type of the recurring task [expected='%s',actual='%s']", line.cmd, result.Type)
		fail = true
	}
	fail = ensureArg(result.Args, line.monday, t)
	fail = ensureArg(result.Args, line.wednesday, t)
	fail = ensureArg(result.Args, line.friday, t)

	if len(result.Args) != 3 {
		t.Errorf("Failed to parse Recurring Task arguments, expected [count=3] got [count=%v]",
			len(result.Args))
		fail = true
	}

	if result.Todo != line.todo {
		t.Errorf("Failed to parse todo line [expected='%s',actual='%s']", line.todo, result.Todo)
	}

	if len(line.attrs) > 0 {
		for key, expectedVal := range line.result {
			actualVal, ok := result.Attr[key]

			if !ok {
				t.Errorf("Expected to have [key=%s] in the attributes, but there is none", key)
				fail = true
				break
			}

			if expectedVal != actualVal {
				t.Errorf("Expected to have [key=%s] with [value=%s], but found [value=%s]", key, expectedVal, actualVal)
				fail = true
				break
			}
		}
	}

	if fail {
		t.FailNow()
	}
}

func ensureArg(args []string, expected string, t *testing.T) bool {
	for _, actual := range args {
		if actual == expected {
			return false // fail == false
		}
	}

	t.Errorf("Failed to parse arguments, expected to find [arg='%s'] in [args=%v] but haven't", expected, args)
	return true // fail == true
}

func TestReturnsErrorIfWrongSyntax(t *testing.T) {
	line := "blah blah blah"

	_, err := parse(line)

	ensureError(err, t)
}

func TestReturnsErrorIfTodoBodyIsEmpty(t *testing.T) {
	line := "blah blah blah:"

	_, err := parse(line)

	ensureError(err, t)
}

func TestReturnsErrorIfLineStartsWithSpace(t *testing.T) {
	line := " blah blah blah:do a thing"

	_, err := parse(line)

	ensureError(err, t)
}

func TestReturnsErrorIfArgumentListEndsWithSpace(t *testing.T) {
	line := "blah blah blah :do a thing"

	_, err := parse(line)

	ensureError(err, t)
}

func ensureError(err error, t *testing.T) {
	if err == nil {
		t.Errorf("Parser accepted syntax which it not supposed to")
		t.Fail()
	}
}
