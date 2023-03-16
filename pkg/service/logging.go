package service

import (
	"bufio"
	"io"
	"log"
	"os"
)

type Logger struct {
	stdoutInput io.ReadCloser
	stderrInput io.ReadCloser
}

func (l *Logger) Bind(stdout, stderr io.ReadCloser) {
	l.stdoutInput = stdout
	l.stderrInput = stderr
}

func (l *Logger) Listen() {
	stdoutChannel := readToChannel(l.stdoutInput)
	stderrChannel := readToChannel(l.stderrInput)

	outLog := log.New(os.Stdout, "STDOUT", 0)
	errLog := log.New(os.Stderr, "STDERR", 0)

	for {
		select {
		case stdout := <-stdoutChannel:
			outLog.Println(stdout)
		case stderr := <-stderrChannel:
			errLog.Println(stderr)
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
