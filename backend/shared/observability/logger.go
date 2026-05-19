package observability

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct {
	service string
	logger  *log.Logger
}

type LogEntry struct {
	Timestamp     time.Time      `json:"timestamp"`
	Service       string         `json:"service"`
	EventID       string         `json:"event_id,omitempty"`
	CorrelationID string         `json:"correlation_id,omitempty"`
	AggregateID   string         `json:"aggregate_id,omitempty"`
	Status        string         `json:"status"`
	Message       string         `json:"message"`
	Fields        map[string]any `json:"fields,omitempty"`
}

func NewLogger(service string) *Logger {
	return &Logger{service: service, logger: log.New(os.Stdout, "", 0)}
}

func (l *Logger) Info(message string, entry LogEntry) {
	entry.Timestamp = time.Now().UTC()
	entry.Service = l.service
	if entry.Status == "" {
		entry.Status = "ok"
	}
	entry.Message = message
	l.write(entry)
}

func (l *Logger) Error(message string, entry LogEntry) {
	entry.Timestamp = time.Now().UTC()
	entry.Service = l.service
	entry.Status = "error"
	entry.Message = message
	l.write(entry)
}

func (l *Logger) write(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		l.logger.Print(`{"status":"error","message":"failed to marshal log"}`)
		return
	}
	l.logger.Print(string(data))
}
