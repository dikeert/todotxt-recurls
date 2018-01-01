package main

import (
	"log"
	"os"
	"path/filepath"
)

var logger *log.Logger

const LOG_FILE = "todotxt-recurls.log"

func init() {
	setupLogging(TodoEnviron())
}

func setupLogging(env *TodoTXT) {
	logpath := filepath.Join(env.DirPath, LOG_FILE)
	if file, err := os.OpenFile(logpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm); err != nil {
		log.Fatal(err)
	} else {
		logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		logger.Printf("Logger ready [path=%s]", logpath)
	}
}
