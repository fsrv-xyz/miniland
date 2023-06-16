package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Event struct {
	Message any `json:"message"`
}

func ServerSendEventsHandlerBuilder(events <-chan Event) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		// Set the necessary headers to allow SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		for event := range events {
			buffer := bytes.NewBufferString("data: ")

			if err := json.NewEncoder(buffer).Encode(event); err != nil {
				log.Println(err)
				return
			}
			buffer.WriteString("\n\n")

			if _, err := buffer.WriteTo(w); err != nil {
				log.Println(err)
				return
			}

			// Flush the data immediately instead of buffering it
			flusher.Flush()
		}
	}
}
