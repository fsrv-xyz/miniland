package service

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const LOG_DIR = "/log/"

type Logger struct {
	stdoutInput io.ReadCloser
	stderrInput io.ReadCloser
}

type LogFmt struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func (logfmt LogFmt) JsonString() string {
	jsonString, _ := json.Marshal(logfmt)
	return string(jsonString)
}

func (l *Logger) Bind(stdout, stderr io.ReadCloser) {
	l.stdoutInput = stdout
	l.stderrInput = stderr
}

func (l *Logger) Listen(ctx context.Context, title string) {

	stdoutChannel := readToChannel(l.stdoutInput)
	stderrChannel := readToChannel(l.stderrInput)

	logFile, logFileOpenError := os.OpenFile(path.Join(LOG_DIR, title+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if logFileOpenError != nil {
		log.Fatal(logFileOpenError)
	}
	outLog := log.New(logFile, "", 0)

	for {
		select {
		case <-ctx.Done():
			logFile.Close()
			return
		case stdout := <-stdoutChannel:
			outLog.Println(LogFmt{Timestamp: time.Now().Format(time.RFC3339), Level: "STDOUT", Message: stdout}.JsonString())
		case stderr := <-stderrChannel:
			outLog.Println(LogFmt{Timestamp: time.Now().Format(time.RFC3339), Level: "STDERR", Message: stderr}.JsonString())
		}
	}
}

func readToChannel(reader io.ReadCloser) chan string {
	channel := make(chan string)
	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			channel <- scanner.Text()
		}
	}()
	return channel
}
