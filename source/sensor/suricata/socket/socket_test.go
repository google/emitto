// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package socket

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"testing"
)

// fakeSocketServer mimiks a Suricata Unix socket server. It handles receiving a version message
// and a subsequent command.
// cmdSize is the expected Command message size for exiting the read loop.
func fakeSocketServer(t *testing.T, conn net.Conn, resp *Response, cmdSize int) {
	// Receive command.
	in := make([]byte, 0, 128)
	for len(in) < cmdSize {
		tmp := make([]byte, 64)
		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		in = append(in, tmp[:n]...)
	}

	// Send response.
	out, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	conn.Write(out)
	conn.Close()
}

func TestVersion(t *testing.T) {
	buf, err := json.Marshal(&Version{versionID})
	if err != nil {
		t.Fatal(err)
	}

	server, client := net.Pipe()
	defer client.Close()
	go fakeSocketServer(t, server, &Response{
		Return:  "OK",
		Message: "1.0",
	}, len(buf))

	s := &Socket{conn: client}
	if err := s.version(); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidCommandName(t *testing.T) {
	_, client := net.Pipe()
	defer client.Close()
	s := &Socket{conn: client}
	_, err := s.Send(&Command{Name: "list-iface"})
	if err == nil {
		t.Fatal(err)
	}
}

func TestInvalidResponse(t *testing.T) {
	cmd := &Command{Name: ReloadRules}
	buf, err := json.Marshal(cmd)
	if err != nil {
		t.Fatal(err)
	}

	server, client := net.Pipe()
	defer client.Close()
	go fakeSocketServer(t, server, &Response{
		Return:  "NOK",
		Message: "error",
	}, len(buf))

	s := &Socket{conn: client}
	if _, err := s.Send(cmd); err == nil {
		t.Fatal(err)
	}
}

func TestConnectTimeout(t *testing.T) {
	p := "/tmp/nonexistent.socket"
	s := New(p)
	if err := s.Connect(); !strings.Contains(err.Error(), "connect: no such file or directory") {
		t.Fatalf("expected connection timeout from: %s", p)
	}
}
