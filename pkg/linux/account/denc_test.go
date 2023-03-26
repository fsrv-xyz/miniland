package account

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStringer struct {
	value string
}

func (m mockStringer) String() string {
	return m.value
}

type mockWriter struct {
	io.Writer
	err error
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.Writer.Write(p)
}

func TestEncoder(t *testing.T) {
	t.Run("normal functionality", func(t *testing.T) {
		var buf bytes.Buffer
		entries := []fmt.Stringer{
			mockStringer{"A"},
			mockStringer{"B"},
			mockStringer{"C"},
		}

		err := encoder(&buf, entries)
		require.NoError(t, err)

		expected := strings.Join([]string{"A", "B", "C"}, "\n") + "\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("handle error", func(t *testing.T) {
		var buf bytes.Buffer
		entries := []fmt.Stringer{
			mockStringer{"A"},
			mockStringer{"B"},
			mockStringer{"C"},
		}

		mockErr := errors.New("write error")
		w := &mockWriter{Writer: &buf, err: mockErr}

		err := encoder(w, entries)
		require.Error(t, err)
		assert.Contains(t, err.Error(), mockErr.Error())
	})

	t.Run("empty slice", func(t *testing.T) {
		var buf bytes.Buffer
		var entries []fmt.Stringer

		err := encoder(&buf, entries)
		require.NoError(t, err)
		assert.Empty(t, buf.String())
	})
}

func TestDecoder(t *testing.T) {
	t.Run("normal functionality with groupEntry", func(t *testing.T) {
		input := strings.Join([]string{"A", "B", "C"}, "\n") + "\n"
		r := strings.NewReader(input)

		lineParser := func(line string) (any, error) {
			return groupEntry{Name: line}, nil
		}

		entries, err := decoder[groupEntry](r, lineParser)
		require.NoError(t, err)

		expected := []groupEntry{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
		}
		assert.Equal(t, expected, entries)
	})

	t.Run("normal functionality with passwordEntry", func(t *testing.T) {
		input := strings.Join([]string{"A", "B", "C"}, "\n") + "\n"
		r := strings.NewReader(input)

		lineParser := func(line string) (any, error) {
			return passwdEntry{Username: line}, nil
		}

		entries, err := decoder[passwdEntry](r, lineParser)
		require.NoError(t, err)

		expected := []passwdEntry{
			{Username: "A"},
			{Username: "B"},
			{Username: "C"},
		}
		assert.Equal(t, expected, entries)
	})

	t.Run("false type", func(t *testing.T) {
		input := strings.Join([]string{"A", "B", "C"}, "\n") + "\n"
		r := strings.NewReader(input)

		lineParser := func(line string) (any, error) {
			return passwdEntry{Username: line}, nil
		}

		_, err := decoder[groupEntry](r, lineParser)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "interface conversion")

	})

	t.Run("handle error", func(t *testing.T) {
		input := strings.Join([]string{"A", "B", "C"}, "\n") + "\n"
		r := strings.NewReader(input)

		mockErr := errors.New("parse error")
		lineParser := func(line string) (any, error) {
			return nil, mockErr
		}

		_, err := decoder[groupEntry](r, lineParser)
		require.Error(t, err)
		assert.Contains(t, err.Error(), mockErr.Error())
	})
}
