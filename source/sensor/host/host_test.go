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

package host

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type fakeConn struct {
	net.Conn
	localAddr *net.UDPAddr
}

func (c fakeConn) LocalAddr() net.Addr { return c.localAddr }
func (fakeConn) Close() error          { return nil }

func fakeHost() *Host {
	return &Host{
		fqdn: "test_1_host_name",
		ip:   net.ParseIP("100.97.26.27"),
	}
}

func TestUpdate(t *testing.T) {
	for _, tt := range []struct {
		desc       string
		osHostname func() (string, error)
		netDial    func(string, string) (net.Conn, error)
		want       *Host
	}{
		{
			desc:       "no change",
			osHostname: func() (string, error) { return "test_1_host_name", nil },
			netDial: func(string, string) (net.Conn, error) {
				updAddr, _ := net.ResolveUDPAddr("udp", "[100.97.26.27]:5688")
				return &fakeConn{localAddr: updAddr}, nil
			},
			want: &Host{
				fqdn: "test_1_host_name",
				ip:   net.ParseIP("100.97.26.27"),
			},
		},
		{
			desc:       "fqdn changed",
			osHostname: func() (string, error) { return "test_2_host_name", nil },
			netDial:    nil,
			want: &Host{
				fqdn: "test_2_host_name",
				ip:   net.ParseIP("100.97.26.27"),
			},
		},
		{
			desc:       "ip changed",
			osHostname: nil,
			netDial: func(string, string) (net.Conn, error) {
				updAddr, _ := net.ResolveUDPAddr("udp", "[10.10.10.10]:5688")
				return &fakeConn{localAddr: updAddr}, nil
			},
			want: &Host{
				fqdn: "test_2_host_name",
				ip:   net.ParseIP("10.10.10.10"),
			},
		},
	} {
		h := fakeHost()

		// Override functions, if applicable.
		if tt.osHostname != nil {
			osHostname = tt.osHostname
		}
		if tt.netDial != nil {
			netDial = tt.netDial
		}

		if err := h.Update(); err != nil {
			t.Errorf("TestUpdate(%s): got unexpected error: %v", tt.desc, err)
		}
		if diff := cmp.Diff(tt.want, h, cmp.AllowUnexported(*tt.want, *h)); diff != "" {
			t.Errorf("TestUpdate(%s): expectation mismatch (-want +got):\n%s", tt.desc, diff)
		}
	}
}

func TestGetters(t *testing.T) {
	h := fakeHost()

	if diff := cmp.Diff("test_1_host_name", h.FQDN()); diff != "" {
		t.Errorf("TestGetters: expectation mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff("100.97.26.27", h.IP().String()); diff != "" {
		t.Errorf("TestGetters: expectation mismatch (-want +got):\n%s", diff)
	}
}
