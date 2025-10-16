package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"testing"
	"time"

	"github.com/quic-go/webtransport-go"
)

func TestWebTransportEcho(t *testing.T) {
	// Create a context that will cancel after the test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		if err := runServer(ctx); err != nil && err != context.Canceled {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait a bit for the server to start
	time.Sleep(100 * time.Millisecond)

	// Create client TLS config (skip verification for self-signed cert)
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create a WebTransport dialer
	dialer := &webtransport.Dialer{
		TLSClientConfig: tlsConf,
	}

	// Dial the WebTransport session
	_, sess, err := dialer.Dial(ctx, "https://localhost:4433/webtransport", nil)
	if err != nil {
		t.Fatalf("Failed to dial WebTransport: %v", err)
	}
	defer sess.CloseWithError(0, "")

	// Open a stream
	stream, err := sess.OpenStream()
	if err != nil {
		t.Fatalf("Failed to open stream: %v", err)
	}
	defer stream.Close()

	// Test data to send
	testData := []byte("Hello, WebTransport!")

	// Send the data
	_, err = stream.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to stream: %v", err)
	}

	// Close the write side (since Stream is Closer, Close closes both, but for echo it should work)
	stream.Close()

	// Read the echoed response
	response := make([]byte, len(testData))
	_, err = io.ReadFull(stream, response)
	if err != nil {
		t.Fatalf("Failed to read from stream: %v", err)
	}

	// Verify the response matches the sent data
	if !bytes.Equal(response, testData) {
		t.Errorf("Expected %q, got %q", testData, response)
	}
}