package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	Message any `json:"message"`
}

func ServerSendEventsHandlerBuilder(events <-chan Event) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(100 * time.Millisecond)
		select {
		case ev := <-events:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.Encode(ev)
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
