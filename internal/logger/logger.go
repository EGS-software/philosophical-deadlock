package logger

import (
	"fmt"
	"time"
)

// ANSI color codes for distinct philosopher output in the terminal
var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
}

const resetColor = "\033[0m"

// Event represents a state change of a philosopher
type Event struct {
	PhilosopherID int
	Action        string
}

// Logger handles thread-safe console output using a channel
type Logger struct {
	eventChan chan Event
}

// New creates and starts a background logger goroutine
func New() *Logger {
	l := &Logger{
		// Buffered channel to avoid blocking the philosophers' execution flow
		eventChan: make(chan Event, 100),
	}
	go l.processEvents()
	return l
}

// Log queues a new state change event to be printed
func (l *Logger) Log(philosopherID int, action string) {
	l.eventChan <- Event{
		PhilosopherID: philosopherID,
		Action:        action,
	}
}

// processEvents runs in a dedicated goroutine, ensuring serialized I/O
func (l *Logger) processEvents() {
	for event := range l.eventChan {
		color := colors[(event.PhilosopherID-1)%len(colors)]
		timestamp := time.Now().Format("15:04:05.000")
		fmt.Printf("%s[%s] [Philosopher %d] %s%s\n", color, timestamp, event.PhilosopherID, event.Action, resetColor)
	}
}
