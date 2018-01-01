package main

import (
	"fmt"
	"os"
)

const CNT = 5
const LINE_CONTENT = "Hello, world!"

var TMP = os.TempDir()

func getRecur() *RecurFile {
	reader, e := OpenRecur(&TodoTXT{DirPath: os.TempDir()})

	if e == nil {
		return reader
	}

	panic(e)
}

type contentWriter func(*os.File)

func newRecurFile(w contentWriter) {
	path := fmt.Sprintf("%s/%s", TMP, RECUR_FILE)

	if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err == nil {
		w(file)
		file.Close()
	} else {
		panic(err)
	}
}

func createRecurFile() {
	newRecurFile(func(w *os.File) {
		for i := 0; i < CNT; i++ {
			w.WriteString(fmt.Sprintln(LINE_CONTENT, i))
		}
	})
}
