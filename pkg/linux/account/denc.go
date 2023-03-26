package account

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

func encoder[T fmt.Stringer](out io.Writer, entries []T) error {
	var encodeError error
	for _, entry := range entries {
		_, err := out.Write([]byte(entry.String() + "\n"))
		if err != nil {
			encodeError = errors.Join(encodeError, err)
		}
	}
	return encodeError
}

type decodable interface {
	passwdEntry | groupEntry
}

type LineParserFunction func(string) (any, error)

func decoder[T decodable](in io.Reader, lineParserFunction LineParserFunction) (entries []T, decodeError error) {
	// catch panics and return them as errors with predictable entries
	defer func() {
		if r := recover(); r != nil {
			entries = nil
			decodeError = fmt.Errorf("panic: %v", r)
		}
	}()

	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		// read a line
		line := scanner.Text()

		// parse the line
		entryPrototype, err := lineParserFunction(line)
		if err != nil {
			decodeError = errors.Join(decodeError, err)
			continue
		}
		// panics here will be caught by the deferred function
		entry := entryPrototype.(T)

		entries = append(entries, entry)
	}

	return entries, decodeError
}
