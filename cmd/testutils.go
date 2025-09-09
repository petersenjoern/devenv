package cmd

import (
	"io"
	"os"
	"strings"
)

// outputCapture encapsulates stdout capture and restoration logic
type outputCapture struct {
	originalStdout *os.File
	reader         *os.File
	writer         *os.File
	done           chan bool
}

// captureOutput starts capturing stdout and returns a handle for restoration
func captureOutput(output *strings.Builder) *outputCapture {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan bool)
	go func() {
		io.Copy(output, r)
		done <- true
	}()

	return &outputCapture{
		originalStdout: originalStdout,
		reader:         r,
		writer:         w,
		done:           done,
	}
}

// restore restores original stdout
func (oc *outputCapture) restore() {
	oc.writer.Close()
	<-oc.done // Wait for the goroutine to finish copying
	os.Stdout = oc.originalStdout
}

// MockCommandExecutor for testing installation without executing real commands
type MockCommandExecutor struct {
	ExecutedCommands []string
	ShouldFail       bool
	FailureError     error
}

func (m *MockCommandExecutor) Execute(command string) error {
	m.ExecutedCommands = append(m.ExecutedCommands, command)
	if m.ShouldFail {
		return m.FailureError
	}
	return nil
}
