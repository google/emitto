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

// Package suricata contains functionality to control and monitor Suricata.
package suricata

import "github.com/google/emitto/source/sensor/suricata/socket"

// Socket represents a Suricata unix socket connection.
type Socket interface {
	// Connect to socket service.
	Connect() error
	// Send command over socket connection.
	Send(*socket.Command) (*socket.Response, error)
	// Close socket connection.
	Close() error
}

// Controller controls Suricata via its Unix socket service.
type Controller struct {
	sock Socket
}

// NewController creates a new Controller.
func NewController(sockAddr string) *Controller {
	return &Controller{socket.New(sockAddr)}
}

// ReloadRules issues a command to Suricata to reload the rules engine with updated rules.
func (c *Controller) ReloadRules() error {
	if err := c.sock.Connect(); err != nil {
		return err
	}
	defer c.sock.Close()

	_, err := c.sock.Send(&socket.Command{
		Name: socket.ReloadRules,
	})
	if err != nil {
		return err
	}

	return nil
}
