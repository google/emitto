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

// Package main initializes and starts the Emitto sensor client.
package main

import (
	"context"
	"flag"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/emitto/source/filestore"
	"github.com/google/emitto/source/sensor/client"

	log "github.com/golang/glog"
)

var (
	// Suricata flags.
	fsSocket       = flag.String("fleetspeak_socket", "", "Fleetspeak client socket")
	suricataSocket = flag.String("suricata_socket", "", "Suricata Unix socket")
	ruleFile       = flag.String("rule_file", "", "Suricata rule file path")
	memoryStorage  = flag.Bool("memory_storage", false, "Use memory store and filestore")

	// Sensor identity flags.
	org  = flag.String("org", "", "Sensor organization")
	zone = flag.String("zone", "", "Sensor zone")

	// Google Cloud Project flags.
	projectID     = flag.String("project_id", "", "Google Cloud project ID")
	storageBucket = flag.String("storage_bucket", "", "Google Cloud Storage bucket for storing rule files")
	credFile      = flag.String("cred_file", "", "Path of the JSON application credential file")

	// Suricata EVE monitoring flags.
	monitorEVE         = flag.Bool("monitor_eve", false, "Monitor EVE logs")
	alertPollingPeriod = flag.Duration("alerts_polling", 10*time.Minute, "Polling interval for checking alerts")
	alertThreshold     = flag.Int("alerts_threshold", 100, "Alerting threshold for Suricata alerts")
	// https://suricata.readthedocs.io/en/suricata-4.1.4/output/eve/eve-json-output.html
	suricataEVELog = flag.String("eve_log", "", "Path of the eve.json file")

	// Heartbeat flags.
	heartbeat              = flag.Bool("heartbeat", false, "Send heartbeat to server")
	heartbeatPollingPeriod = flag.Duration("heartbeat_polling", 10*time.Minute, "Polling interval for sending heartbeats")
)

func main() {
	ctx := context.Background()

	fs, closeFStore := mustGetFileStore(ctx)
	defer closeFStore()
	sc, err := client.New(ctx, *fsSocket, *org, *zone, *ruleFile, *suricataSocket, fs)
	if err != nil {
		log.Exitf("failed to create sensor client: %v", err)
	}

	done := make(chan struct{})
	go func() {
		sc.FSClient.Receive(done)
	}()
	defer close(done)

	if *monitorEVE {
		go monitorEVELog(sc)
	}

	if *heartbeat {
		go heartbeatPolling(sc)
	}

	for msg := range sc.FSClient.Messages() {
		if err := sc.ProcessMessage(ctx, msg); err != nil {
			log.Error(err.Error())
		}
	}
}

func mustGetFileStore(ctx context.Context) (filestore.FileStore, func() error) {
	if *memoryStorage {
		return filestore.NewMemoryFileStore(), func() error { return nil }
	}
	c, err := filestore.NewGCSClient(ctx, *credFile, []string{storage.ScopeReadOnly})
	if err != nil {
		log.Exitf("failed to create Google Cloud Storage client: %v", err)
	}
	return filestore.NewGCSFileStore(*storageBucket, c), c.Close
}

func monitorEVELog(sc *client.Client) {
	for range time.Tick(*alertPollingPeriod) {
		sc.MonitorSurcataEVELog(*alertPollingPeriod, *alertThreshold, *suricataEVELog)
	}
}

func heartbeatPolling(sc *client.Client) {
	for range time.Tick(*heartbeatPollingPeriod) {
		sc.SendHeartbeat()
	}
}
