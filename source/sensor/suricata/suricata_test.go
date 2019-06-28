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

package suricata

import (
	"net"
	"testing"

	"github.com/google/emitto/source/sensor/suricata/socket"
)

type fakeSocket struct {
	conn net.Conn

	sendResponse *socket.Response
}

func (s *fakeSocket) Connect() error { return nil }
func (s *fakeSocket) Close() error   { return nil }

func (s *fakeSocket) Send(*socket.Command) (*socket.Response, error) {
	return s.sendResponse, nil
}

func TestReloadRules(t *testing.T) {
	fs := new(fakeSocket)
	fs.sendResponse = &socket.Response{
		Return: "OK",
	}

	ctrl := &Controller{fs}
	if err := ctrl.ReloadRules(); err != nil {
		t.Fatal(err)
	}
}
