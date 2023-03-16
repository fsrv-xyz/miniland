package service

import (
	"bufio"
	"io"
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
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

	loggeri := zlog.Output(zerolog.New(os.Stdout).With().Timestamp().Logger())

	for {
		select {
		case stdout := <-stdoutChannel:
			loggeri.Info().Msg(stdout)
		case stderr := <-stderrChannel:
			loggeri.Warn().Msg(stderr)
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
