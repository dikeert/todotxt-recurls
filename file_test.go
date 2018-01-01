package main

import (
	"fmt"
	"testing"
)

const REPLACE_CONTENT = "World, Hello!"

func TestIsAbleToUpdateCertainLines(t *testing.T) {
	createRecurFile()

	line := uint(3)
	records := []*RecurRecord{
		{line, fmt.Sprintf("%s %v", REPLACE_CONTENT, line)},
	}
	recur := getRecur()
	recur.Write(records)

	for i, scanner := 0, getRecur().Scanner(); scanner.Scan(); i++ {
		if uint(i) == line {
			expected := fmt.Sprintf("%s %v", REPLACE_CONTENT, i)
			actual := scanner.Text()

			if actual == expected {
				return
			} else {
				t.Errorf("Unexpected record content [expect='%s',actual='%s']", expected, actual)
				t.FailNow()
			}
		}
	}

	t.Fatal("Can't detect any repaced lines")
	t.FailNow()
}

func TestIsAbleToUpdateFirstLine(t *testing.T) {
	createRecurFile()

	line := uint(0)
	records := []*RecurRecord{
		{line, fmt.Sprintf("%s %v", REPLACE_CONTENT, line)},
	}
	recur := getRecur()
	recur.Write(records)

	for i, scanner := 0, getRecur().Scanner(); scanner.Scan(); i++ {
		if uint(i) == line {
			expected := fmt.Sprintf("%s %v", REPLACE_CONTENT, i)
			actual := scanner.Text()

			if actual == expected {
				return
			} else {
				t.Errorf("Unexpected record content [expect='%s',actual='%s']", expected, actual)
				t.FailNow()
			}
		}
	}

	t.Fatal("Can't detect any repaced lines")
	t.FailNow()
}
