package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type TodoTXT struct {
	DirPath    string
	FilePath   string
	ConfigPath string
}

func (me *TodoTXT) AddTodo(todo string) error {
	file, err := os.OpenFile(me.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModeAppend)

	if err != nil {
		return err
	}

	if _, err := file.WriteString(fmt.Sprintf("%s\n", todo)); err != nil {
		return err
	}

	return nil
}

func (me *TodoTXT) String() string {
	return fmt.Sprintf(
		"TodoTXT{DirPath=%s, FilePath=%s, ConfigPath=%s}",
		me.DirPath, me.FilePath, me.ConfigPath)
}

func todoRawEnviron() map[string]string {
	environ := make(map[string]string)

	for _, env := range os.Environ() {
		if strings.Index(env, "TODO") == 0 {
			var key string
			var val bytes.Buffer

			for idx, token := range strings.Split(env, "=") {
				if idx == 0 {
					key = token
				} else {
					if _, err := val.WriteString(token); err != nil {
						panic(err)
					}
				}
			}
			environ[key] = val.String()
		}
	}

	return environ
}

func TodoEnviron() *TodoTXT {
	var env TodoTXT

	for key, val := range todoRawEnviron() {
		switch key {
		case "TODO_DIR":
			env.DirPath = val
		case "TODO_FILE":
			env.FilePath = val
		case "TODOTXT_CFG_FILE":
			env.ConfigPath = val
		}
	}

	return &env
}
