package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	os.Stdout = writer

	outputC := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		outputC <- buf.String()
	}()

	fn()

	_ = writer.Close()
	os.Stdout = originalStdout
	output := <-outputC
	_ = reader.Close()

	return output
}
