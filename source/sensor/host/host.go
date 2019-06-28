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

// Package host provides functionality to obtain basic host information for a host.
package host

import (
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/golang/glog"
)

var (
	// To facilitate unit testing.
	netDial    = net.Dial
	osHostname = os.Hostname
)

// Host contains basic host information.
type Host struct {
	fqdn string
	ip   net.IP
}

// FQDN returns the Host FQDN.
func (h *Host) FQDN() string {
	return h.fqdn
}

// IP returns the Host IP address.
func (h *Host) IP() net.IP {
	return h.ip
}

// New performs an initial update and returns a Host.
func New() (*Host, error) {
	h := new(Host)

	// Wait indefinitely for the host to come online.
	for {
		if err := h.Update(); err != nil {
			log.Errorf("host update failed: %v; retrying after 5 seconds", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return h, nil
}

// Update updates the host.
func (h *Host) Update() error {
	if err := h.updateFQDN(); err != nil {
		return err
	}
	if err := h.updateIP(); err != nil {
		return err
	}
	log.Infof("Host info: hostname: %s, ip: %s", h.fqdn, h.ip)
	return nil
}

// updateFQDN updates the host fqdn.
func (h *Host) updateFQDN() error {
	fqdn, err := osHostname()
	if err != nil {
		return fmt.Errorf("failed to retrieve fqdn: %v", err)
	}
	h.fqdn = fqdn
	return nil
}

// updateIP updates the host IP address.
func (h *Host) updateIP() error {
	conn, err := netDial("udp", "8.8.8.8:9") // RFC863.
	if err != nil {
		return fmt.Errorf("failed to retrieve ip address: %v", err)
	}
	defer conn.Close()
	h.ip = conn.LocalAddr().(*net.UDPAddr).IP
	return nil
}
