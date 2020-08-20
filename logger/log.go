package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/go-errors/errors"
)

// Log is the base structure for logging
type Log struct {
	logFile string
	ch      chan string
	wg      *sync.WaitGroup
}

// NewLogger is the default constructor to create logger
func NewLogger(projectName string) *Log {
	if strings.TrimSpace(projectName) == "" {
		panic("The project name has to be set in arguments.")
	}

	logger := Log{
		logFile: getFileName(projectName),
		ch:      make(chan string, 5),
		wg:      &sync.WaitGroup{},
	}
	go logger.write()
	return &logger
}

// GetFileName returns a filename where logger writes info
func (l Log) GetFileName() string {
	return l.logFile
}

func getFileName(projectName string) string {
	return fmt.Sprintf("log/%s%s.log",
		projectName, time.Now().Format("2006-01-02"))
}

func (l Log) write() {
	err := os.MkdirAll("log", 0777)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(l.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(getErrorMessage(err))
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	for message := range l.ch {
		write(message, writer, l.wg)
	}
}

func write(message string, writer io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(writer, strings.NewReader(message))
	if err != nil {
		panic(getErrorMessage(err))
	}
}

// Close writing messages buffer
func (l Log) Close() {
	l.wg.Wait()
	close(l.ch)
}

// Info writes information messages
func (l Log) Info(message ...interface{}) {
	l.log(Info, message...)
}

// Error writes error messages
func (l Log) Error(err error, message ...interface{}) {
	if len(message) > 0 {
		l.log(Error, message...)
	}
	l.log(Error, errors.Wrap(err, 1).ErrorStack())
}

// LogIfError writes error message if err not nil
func (l Log) LogIfError(err error, message ...interface{}) {
	if err != nil {
		if len(message) > 0 {
			l.log(Error, message...)
		}
		l.log(Error, errors.Wrap(err, 1).ErrorStack())
	}
}

// Fatal writes fatal messages
func (l Log) Fatal(err error, message ...interface{}) {
	if len(message) > 0 {
		l.log(Fatal, message...)
	}
	l.log(Fatal, errors.Wrap(err, 1).ErrorStack())
}

func (l Log) log(level Level, message ...interface{}) {
	l.wg.Add(1)
	l.ch <- getMessage(level, message...)
}

func getErrorMessage(err error) string {
	return getMessage(Error, err)
}

func getMessage(level Level, message ...interface{}) string {
	currentTime := time.Now().Format(time.RFC3339)
	str := fmt.Sprint(message...)
	return fmt.Sprintf("%s [%s]: %s\n", currentTime, level, str)
}
