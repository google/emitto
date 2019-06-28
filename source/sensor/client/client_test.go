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

package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/emitto/source/sensor/host"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/google/emitto/source/sensor/proto"
)

type fakeSuricataController struct{}

func (s *fakeSuricataController) ReloadRules() error { return nil }

func TestSocketReloadRules(t *testing.T) {
	c := &Client{
		ctrl: &fakeSuricataController{},
	}

	got := c.reloadRules()
	want := status.New(codes.OK, "OK")
	if diff := cmp.Diff(want, got, cmp.Comparer(proto.Equal), cmp.AllowUnexported(*want, *got)); diff != "" {
		t.Fatalf("expectation mismatch:\n%s", diff)
	}
}

func TestParseLogLine(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		line    string
		want    string
		wantErr bool
	}{
		{
			desc: "well-formed JSON",
			line: `{"timestamp": "2019-05-13T14:12:19.384640+0000", "event_type": "alert"}`,
			want: `timestamp:"2019-05-13T14:12:19.384640+0000" event_type:"alert" `,
		}, {
			desc:    "malformed JSON; missing parenthesis",
			line:    `{"dest_port":234234`,
			wantErr: true,
		}, {
			desc:    "malformed JSON",
			line:    `<Warning> -- [ERRCODE: SC_ERR_EVENT_ENGINE(210)]`,
			wantErr: true,
		},
	} {
		eve, err := parseLogLine(tt.line)

		if (err != nil) != tt.wantErr {
			t.Errorf("%s: got err=%v, wantErr=%t", tt.desc, err, tt.wantErr)
		}
		if err != nil {
			continue
		}
		if diff := cmp.Diff(tt.want, eve.String()); diff != "" {
			t.Errorf("parseLogLine(%v) expectation mismatch (+want -got)\n%s", tt.line, diff)
		}
	}
}

// Creates temporary log file with 10 logs 1 second apart, checks if the sensor emits the alert.
func TestMonitorSurcataEVELog(t *testing.T) {

	f := &fakeFleetspeakClient{}
	client := &Client{
		FSClient: f,
		host:     &host.Host{},
	}

	d, err := ioutil.TempDir("/tmp", "eve-logs")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(d)

	var lines []string
	now := time.Now().UTC()
	for i := 0; i < 10; i++ {
		timestamp := now.Add(time.Duration(-i) * time.Second).Format("2006-01-02T15:04:05.999999+0000")
		lines = append(lines, fmt.Sprintf(`{"timestamp": %q, "event_type": "alert"}`, timestamp))
	}

	logFile := filepath.Join(d, "eve.json")
	if err := ioutil.WriteFile(logFile, []byte(strings.Join(lines, "\n")), 0666); err != nil {
		t.Fatalf("failed to write in file %v: %v", logFile, err)
	}
	client.MonitorSurcataEVELog(time.Minute, 9, logFile)
	if len(f.Msgs) != 1 {
		t.Errorf("TestMonitorSurcataEVELog(%v, %v, %v), expected to emit alert message", time.Minute, 9, logFile)
	}
	f.Msgs = f.Msgs[:0]
	client.MonitorSurcataEVELog(time.Second, 1, logFile)
	if len(f.Msgs) != 0 {
		t.Errorf("TestMonitorSurcataEVELog(%v, %v, %v), emitted unexpected alert messages: %+v", time.Minute, 1, logFile, f.Msgs)
	}
}

// Creates temporary file, calls createBackup and compares content of returned file to original file.
func TestCreateBackup(t *testing.T) {
	d, err := ioutil.TempDir("/tmp", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(d)

	srcData := []byte{1, 2, 3}
	srcFile := filepath.Join(d, "src.txt")
	if err := ioutil.WriteFile(srcFile, srcData, 0644); err != nil {
		t.Fatalf("failed to write in file %v: %v", srcFile, err)
	}
	backup, err := createBackup(srcFile)
	if err != nil {
		t.Errorf("createBackup(%v) returned an error: %v", srcFile, err)
	}
	backupData, err := ioutil.ReadFile(backup)
	if err != nil {
		t.Errorf("error reading file %v: %v", srcFile, err)
	}
	if diff := cmp.Diff(srcData, backupData); diff != "" {
		t.Errorf("backup file does not match source (-want +got):\n%s", diff)
	}
}

type fakeFleetspeakClient struct {
	FleetspeakClient
	Msgs []*pb.SensorMessage
}

func (c *fakeFleetspeakClient) SendMessage(m *pb.SensorMessage) (string, error) {
	c.Msgs = append(c.Msgs, m)
	return "ok", nil
}
