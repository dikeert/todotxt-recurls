package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const RECUR_FILE = "recur.txt"

type RecurRecord struct {
	ID   uint
	Line string
}

type RecurFile struct {
	dir  string
	file *os.File
	env  *TodoTXT
}

func (me *RecurFile) Scanner() *bufio.Scanner {
	me.file.Seek(0, 0)
	return bufio.NewScanner(me.file)
}

func (me *RecurFile) Writer() *bufio.Writer {
	return bufio.NewWriter(me.file)
}

func (me *RecurFile) ReadAll(target chan *RecurRecord) {
	scanner := me.Scanner()
	for count := uint(0); scanner.Scan(); count++ {
		target <- &RecurRecord{count, scanner.Text()}
	}

	close(target)
}

func (me *RecurFile) WriteAll(source chan *RecurRecord) {
	var records []*RecurRecord

	for record := range source {
		records = append(records, record)
	}

	if len(records) > 0 {
		logger.Printf("Ready to update [count=%v] records", len(records))
		me.Write(records)
	} else {
		logger.Printf("No Recurring Records to update in the file")
	}
}

func (me *RecurFile) Write(records []*RecurRecord) {
	const TEMP_FILE = ".recur.txt.tmp"

	tmp := filepath.Join(me.dir, TEMP_FILE)
	var writer *bufio.Writer

	if file, err := os.OpenFile(tmp, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm); err == nil {
		writer = bufio.NewWriter(file)
	} else {
		panic(err)
	}

	for line, scanner := uint(0), me.Scanner(); scanner.Scan(); line++ {
		text := scanner.Text()

		for _, record := range records {
			if record.ID == line {
				logger.Printf(
					"Found line [ID=%v, value=%s] to replace current line [ID=%v, value=%s]",
					record.ID, record.Line, line, text)
				text = record.Line
				break
			}
		}

		writer.WriteString(fmt.Sprintf("%s\n", text))
	}
	writer.Flush()

	me.file.Close()
	me.file = nil

	os.Rename(tmp, filepath.Join(me.dir, RECUR_FILE))
	if file, err := doOpenRecur(path.Join(me.dir, RECUR_FILE)); err == nil {
		me.file = file
	} else {
		panic(err)
	}
}

func doOpenRecur(path string) (*os.File, error) {
	if file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open [file=%s]: %s", path, err))
	} else {
		return file, nil
	}
}

func OpenRecur(env *TodoTXT) (*RecurFile, error) {
	path := filepath.Join(env.DirPath, RECUR_FILE)

	var file *os.File
	var err error

	if file, err = doOpenRecur(path); err == nil {
		return &RecurFile{env.DirPath, file, env}, nil
	} else {
		return nil, err
	}
}
