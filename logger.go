package restart

import (
	"io/ioutil"
	"log"
	"os"
)

var logger *log.Logger
var logFile *os.File

func init() {
	logger = log.New(os.Stderr, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
}

func NewLogger(fileName string) error {
	if logFile = os.NewFile(4, ""); os.Getenv("graceful") == "true" && logFile != nil {
		logger = log.New(logFile, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
		return nil
	}
	if fileName == "" {
		logger = log.New(ioutil.Discard, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
		return nil
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	logger = log.New(file, "[Restart]", log.Llongfile|log.Lmicroseconds|log.Ldate)
	logFile = file
	return nil
}
