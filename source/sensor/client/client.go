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

// Package client contains sensor client functionality.
package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/emitto/source/filestore"
	"github.com/google/emitto/source/sensor/fleetspeak"
	"github.com/google/emitto/source/sensor/host"
	"github.com/google/emitto/source/sensor/suricata"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/golang/glog"
	pb "github.com/google/emitto/source/sensor/proto"
	evepb "github.com/google/emitto/source/sensor/suricata/proto"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
)

// SuricataController represents a Suricata controller.
type SuricataController interface {
	// ReloadRules reloads Suricata rules.
	ReloadRules() error
}

// FleetspeakClient represents a Fleetspeak client.
type FleetspeakClient interface {
	// SendMessage sends a message to Fleetspeak.
	SendMessage(m *pb.SensorMessage) (string, error)
	// Receive initiates receiving messages from Fleetspeak.
	Receive(done <-chan struct{})
	// Messages provides access to the messages received from Fleetspeak.
	Messages() chan *fspb.Message
}

// Client represents a Emitto sensor client.
type Client struct {
	FSClient  FleetspeakClient
	ctrl      SuricataController
	ruleStore filestore.FileStore
	host      *host.Host
	org       string
	zone      string
	ruleFile  string
}

// New creates a new Emitto sensor client.
func New(ctx context.Context, fleetspeakSocket, org, zone, ruleFile, suricataSocket string, filestore filestore.FileStore) (*Client, error) {
	h, err := host.New()
	if err != nil {
		return nil, fmt.Errorf("failed to created new Host: %v", err)
	}
	return &Client{
		FSClient:  fleetspeak.New(fleetspeakSocket),
		ctrl:      suricata.NewController(suricataSocket),
		ruleStore: filestore,
		host:      h,
		org:       org,
		zone:      zone,
		ruleFile:  ruleFile,
	}, nil
}

// ProcessMessage handles a Fleetspeak message from Emitto.
func (c *Client) ProcessMessage(ctx context.Context, m *fspb.Message) error {
	var req pb.SensorRequest
	if err := ptypes.UnmarshalAny(m.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal Fleetspeak message: %v", err)
	}
	switch t := req.Type.(type) {
	case *pb.SensorRequest_DeployRules:
		log.Infof("Received DeployRules request %q", req.GetId())
		return c.sendResponse(req.GetId(), c.deployRules(ctx, t.DeployRules.GetRuleFile()))
	case *pb.SensorRequest_ReloadRules:
		log.Infof("Received ReloadRules request %q", req.GetId())
		return c.sendResponse(req.GetId(), c.reloadRules())
	default:
		return c.sendResponse(req.GetId(), status.New(codes.InvalidArgument, fmt.Sprintf("unknown request type: %T", t)))
	}
}

// sendResponse sends a SensorResponse to the Fleetspeak client.
func (c *Client) sendResponse(id string, s *status.Status) error {
	if err := c.host.Update(); err != nil {
		return fmt.Errorf("failed to update host info: %v", err)
	}
	resp := &pb.SensorResponse{
		Id:     id,
		Time:   ptypes.TimestampNow(),
		Host:   c.getHostInfo(),
		Status: s.Proto(),
	}
	msg := &pb.SensorMessage{
		Id: uuid.New().String(),
		Type: &pb.SensorMessage_Response{
			Response: resp,
		},
	}
	_, err := c.FSClient.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send response: %+v", resp)
	}
	log.V(1).Infof("Sent response: %+v", resp)
	return nil
}

// deployRules fetches an updated rule file, replaces the existing rule file with the updated
// version, and then issues a command for Suricata to reload the rule engine.
func (c *Client) deployRules(ctx context.Context, ruleFile string) *status.Status {
	rules, err := c.ruleStore.GetRuleFile(ctx, ruleFile)
	if err != nil {
		return status.New(codes.NotFound, fmt.Sprintf("failed to download rules from Cloud Storage: %v", err))
	}
	if _, err := os.Stat(c.ruleFile); os.IsNotExist(err) {
		return status.New(codes.NotFound, fmt.Sprintf("rule file does not exist %q", c.ruleFile))
	}
	// Backup existing rule before writing.
	backup, err := createBackup(c.ruleFile)
	if err != nil {
		return status.New(codes.Internal, fmt.Sprintf("failed to create backup for the rule file %q: %v", c.ruleFile, err))
	}

	if err := ioutil.WriteFile(c.ruleFile, rules, 0644); err != nil {
		// Restore rule from backup file.
		if err := copyFile(backup, c.ruleFile); err != nil {
			log.Errorf("failed to restore rule file %q from backup: %q", c.ruleFile, backup)
		}
		return status.New(codes.Internal, fmt.Sprintf("failed to write rule file to disk: %v", err))
	}
	log.Infof("Successfully wrote new rules to %q", c.ruleFile)
	if err := c.reloadRules(); err != nil {
		return status.New(codes.FailedPrecondition, fmt.Sprintf("failed to reload Suricata rules: %v", err))
	}
	log.Info("Successfully reloaded Suricata rules")
	return status.New(codes.OK, "OK")
}

// reloadRules reloads rules via the Suricata socket.
func (c *Client) reloadRules() *status.Status {
	if err := c.ctrl.ReloadRules(); err != nil {
		return status.New(codes.FailedPrecondition, fmt.Sprintf("failed to issue socket command: %v", err))
	}
	return status.New(codes.OK, "OK")
}

func (c *Client) getHostInfo() *pb.Host {
	return &pb.Host{
		Fqdn: c.host.FQDN(),
		Ip:   c.host.IP().String(),
		Org:  c.org,
		Zone: c.zone,
	}
}

// MonitorSurcataEVELog reads logs from last `d` duration of time,
// and sends a notification alert if the threshold is exceeded
func (c *Client) MonitorSurcataEVELog(d time.Duration, threshold int, logFile string) {
	lines, err := readLines(logFile)
	if err != nil {
		log.Errorf("Error reading lines from log file %v: %v", logFile, err)
		return
	}

	now := time.Now().UTC()
	alerts := 0
	for i := len(lines) - 1; i >= 0; i-- {
		if alerts > threshold {
			break
		}
		eve, err := parseLogLine(lines[i])
		if err != nil {
			log.Error(err)
			c.sendAlertNotification(err.Error())
			return
		}

		// https://suricata.readthedocs.io/en/suricata-4.1.4/output/eve/eve-json-format.html#event-type-alert
		if eve.EventType != "alert" {
			continue
		}

		t, err := time.Parse("2006-01-02T15:04:05.999999+0000", eve.Timestamp)
		if err != nil {
			e := fmt.Errorf("Error parsing timestamp %q: %v", eve.Timestamp, err)
			log.Error(e)
			c.sendAlertNotification(e.Error())
		}

		// Only check last logs within specified duration.
		if now.Sub(t) > d {
			break
		}

		alerts++
	}
	// If the threshold is exceeded, send notification.
	if alerts > threshold {
		c.sendAlertNotification(fmt.Sprintf("Sensor has detected: %v alerts in last: %v", alerts, d))
	}
}

// SendHeartbeat sends a heartbeat message to the server.
func (c *Client) SendHeartbeat() {
	c.FSClient.SendMessage(&pb.SensorMessage{
		Id: uuid.New().String(),
		Type: &pb.SensorMessage_Heartbeat{
			Heartbeat: &pb.Heartbeat{
				Time: ptypes.TimestampNow(),
				Host: c.getHostInfo(),
			},
		},
	})
}

func (c *Client) sendAlertNotification(alert string) {
	c.FSClient.SendMessage(&pb.SensorMessage{
		Id: uuid.New().String(),
		Type: &pb.SensorMessage_Alert{
			Alert: &pb.SensorAlert{
				Time:   ptypes.TimestampNow(),
				Host:   c.getHostInfo(),
				Status: status.New(codes.Internal, alert).Proto(),
			},
		},
	})
}

// createBackup creates copy of filename to filename_YYMMDD_number, with number randomly
// chosen such that the file name is unique and returns the chosen file name.
func createBackup(filename string) (string, error) {
	// create backup file.
	f, err := ioutil.TempFile(filepath.Dir(filename), fmt.Sprintf("%v_%v_", filepath.Base(filename), time.Now().Format("060102")))
	if err != nil {
		return "", err
	}
	backup := f.Name()
	f.Close()
	if err := copyFile(filename, backup); err != nil {
		return "", err
	}
	return backup, nil
}

func copyFile(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, data, 0644)
}

// Note: currently suricata eve.json file is ~1.5GB, reading ~1GB file with 400k lines
// takes less than a second, so this should not become performance bottleneck in near future.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// parseLogLine unmarshals json and returns as EVE protobuf message.
func parseLogLine(line string) (*evepb.EVE, error) {
	var eve evepb.EVE
	if err := json.Unmarshal([]byte(line), &eve); err != nil {
		return nil, fmt.Errorf("cannot unmarshal json string %q: %v", line, err)
	}
	return &eve, nil
}
