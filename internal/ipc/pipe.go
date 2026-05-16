//go:build windows

package ipc

import (
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

// PipeName is the named pipe address used for single-instance forwarding.
const PipeName = `\\.\pipe\SwitchyPipe`

// Message is the JSON payload sent over the named pipe.
type Message struct {
	URL string `json:"url"`
}

// TryForward attempts to forward a URL to an existing Switchy instance via named pipe.
func TryForward(rawURL string) bool {
	conn, err := winio.DialPipe(PipeName, timeoutPtr(2*time.Second))
	if err != nil {
		return false
	}
	defer conn.Close()

	msg := Message{URL: rawURL}
	_ = json.NewEncoder(conn).Encode(msg)
	return true
}

// ListenAndServe listens on the named pipe and dispatches URLs to the handler.
func ListenAndServe(handler func(rawURL string)) error {
	cfg := &winio.PipeConfig{SecurityDescriptor: "D:P(A;;GA;;;WD)"}
	ln, err := winio.ListenPipe(PipeName, cfg)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, handler)
	}
}

func handleConn(conn net.Conn, handler func(string)) {
	defer conn.Close()
	data, err := io.ReadAll(conn)
	if err != nil {
		return
	}
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	if msg.URL != "" {
		handler(msg.URL)
	}
}

func timeoutPtr(d time.Duration) *time.Duration { return &d }
