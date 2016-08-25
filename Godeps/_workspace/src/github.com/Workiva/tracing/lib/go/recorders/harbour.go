package recorders

import (
	"encoding/json"
	"github.com/Workiva/tracing/lib/go"
	"log"
	"net"
	"os"
	"time"
)

const (
	tracingSocket = "/dev/harbour/tracing"
	RFC3339Micro  = "2006-01-02T15:04:05.999999Z07:00"
)

type HarbourRecorder struct {
	socketPath  string
	messages    chan []byte
}

func NewHarbourRecorder(socketPath string) *HarbourRecorder {
	if socketPath == "" {
		socketPath = tracingSocket
	}
	recorder := &HarbourRecorder{
		socketPath:  socketPath,
		messages:    make(chan []byte),
	}
	go recorder.sendMessages()
	return recorder
}

func (h *HarbourRecorder) sendMessages() {
	nextTry := time.After(time.Nanosecond)
	toStderr := true
	retries := uint(0)
	var conn net.Conn

	for {
		select {
		case msg := <- h.messages:
			if toStderr == false {
				if _, err := conn.Write(msg); err != nil {
					toStderr = true
					nextTry = time.After(time.Second)
				}
			}
			if toStderr == true {
				os.Stderr.WriteString("Failed to write span: ")
				os.Stderr.Write(msg)
				os.Stderr.WriteString("\n")
			}
		case <- nextTry:
			c, err := h.connect()
			if err == nil && c != nil {
				conn = c
				toStderr = false
			} else {
				retries = retries + 1
				ival := 1 << retries
				if ival > 32 {
					ival = 32
				}
				nextTry = time.After(time.Duration(ival) * time.Second)
			}
		}
	}
	return
}

func (h *HarbourRecorder) connect() (net.Conn, error) {
	// log.Printf("Try to connect to <%s>", h.socketPath)

	if _, err := os.Stat(h.socketPath); err != nil {
		// log.Printf("Couldn't stat <%s>: %s", h.socketPath, err)
		return nil, err
	}
	conn, err := net.DialTimeout("unixgram", h.socketPath, time.Second)
	if err != nil {
		// log.Printf("Couldn't dial <%s>: %s", h.socketPath, err)
		return nil, err
	}
	// log.Printf("Connected to <%s> successfully", h.socketPath)
	return conn, nil
}

func (h *HarbourRecorder) RecordSpan(span tracing.RawSpan) {
	span.Tags["end"] = span.Start.Add(span.Duration).Format("2006-01-02T15:04:05.000000000Z")

	out := map[string]interface{}{
		"type": "tracing",
		"span": tracing.NewJSONSpan(span),
	}
	blob, err := json.Marshal(out)
	if err != nil {
		log.Printf("Failed to serialize span to json: %s", err)
		log.Printf("Span was: %+v", span)
		return
	}
	h.messages <- blob
	return
}
