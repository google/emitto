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

// Package socket contains functionality to send commands to Suricata via its Unix socket.
//
// Proper usage of the socket:
//   1. Connect()
//   2. Send()
//   3. Close()
package socket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	log "github.com/golang/glog"
)

var (
	retryAttempts = 3
	retryInterval = 5 * time.Second
)

const versionID = "1.0"

// Version represents a version message, which must be sent and return "OK" before sending commands.
type Version struct {
	ID string `json:"version,"`
}

// CommandName represents a Suricata Unix socket command name.
//
// https://suricata.readthedocs.io/en/suricata-4.0.5/unix-socket.html
type CommandName string

// Suricata socket commands.
const (
	ReloadRules CommandName = "reload-rules"
)

// validCommands contains the currently supported socket commands.
var validCommands = map[CommandName]bool{
	ReloadRules: true,
}

// Command represents a Suricata Unix socket command.
//
// Protocol: https://redmine.openinfosecfoundation.org/projects/suricata/wiki/Unix_Socket#Protocol
type Command struct {
	Name CommandName       `json:"command,"`
	Args map[string]string `json:"arguments,omitempty"`
}

// Response represents a Suricata Unix socket command response.
type Response struct {
	Return  string `json:"return,"`
	Message string `json:"message,string"`
}

// Socket represents a Suricata Unix socket server connection.
type Socket struct {
	addr string
	conn net.Conn
}

// New creates a new Socket.
func New(addr string) *Socket {
	return &Socket{addr: addr}
}

// Connect dials the Suricata Unix socket and prepares the connection for receiving commands.
func (s *Socket) Connect() error {
	if err := retry(retryAttempts, retryInterval, func() error {
		c, err := net.Dial("unix", s.addr)
		if err != nil {
			return fmt.Errorf("failed to connect to Suricata socket (%s): %v", s.addr, err)
		}
		s.conn = c
		return nil
	}); err != nil {
		return err
	}

	return s.version()
}

// Close closes the Suricata Unix socket connection.
func (s *Socket) Close() error {
	return s.conn.Close()
}

// version sends a Suricata version message to establish the communication protocol.
// The protocol is defined as the following:
//  1. Client connects to the socket.
//  2. Client sends a version message: { "version": "$VERSION_ID" }.
//  3. Server answers with { "return": "OK|NOK" }.
func (s *Socket) version() error {
	buf, err := json.Marshal(&Version{versionID})
	if err != nil {
		return err
	}
	if _, err := s.conn.Write(buf); err != nil {
		return err
	}

	resp, err := ioutil.ReadAll(s.conn)
	if err != nil {
		return err
	}
	r := new(Response)
	if err := json.Unmarshal(resp, r); err != nil {
		return err
	}
	if r.Return != "OK" {
		return fmt.Errorf("failed to establish communication protocol with Suricata: %s: %s", r.Return, r.Message)
	}
	return nil
}

// Send sends a command to Suricata and returns its response.
func (s *Socket) Send(cmd *Command) (*Response, error) {
	if _, ok := validCommands[cmd.Name]; !ok {
		return nil, fmt.Errorf("unsupported command type: %s", cmd.Name)
	}

	buf, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	if _, err := s.conn.Write(buf); err != nil {
		return nil, err
	}

	// The read blocks until it receives an EOF or returns an error.
	resp, err := ioutil.ReadAll(s.conn)
	if err != nil {
		return nil, err
	}
	r := new(Response)
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}
	if r.Return != "OK" {
		return nil, fmt.Errorf("received an error response from command (%+v): %s: %s", cmd, r.Return, r.Message)
	}

	return r, nil
}

// retry calls a function up to the specified number of attempts if the call encounters an error.
// Each retry will sleep for a specified duration.
func retry(attempts int, sleep time.Duration, do func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = do(); err == nil {
			return nil
		}
		time.Sleep(sleep)
		log.Infof("retrying after error: %v", err)
	}
	return fmt.Errorf("failed after %d attempts; last error: %v", attempts, err)
}
