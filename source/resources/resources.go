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

// Package resources contains common objects and conversion functions.
package resources

const (
	// fleetspeakPrefix is the default label prefix prepended to all client labels.
	fleetspeakPrefix = "alphabet-"
	// LocationNamePrefix is the Fleetspeak label prefix for sensor location name.
	LocationNamePrefix = fleetspeakPrefix + "location-name-"
	// LocationZonePrefix is the Fleetspeak label prefix for sensor location zone.
	LocationZonePrefix = fleetspeakPrefix + "location-zone-"
)

// Location defines an arbirary organization of sensors, segmented into a least one zone.
type Location struct {
	// The unique name of the location, e.g. "company1".
	Name string `mutable:"false"`
	// The list of zones or "segments" to organize sensors, e.g. {"dmz", "prod"}.
	Zones []string `mutable:"true"`
	// Last modified time of the message. Applied by the Store.
	LastModified string `mutable:"true"`
}

// ZoneFilterMode defines how the location zones will be selected.
type ZoneFilterMode string

const (
	// All is to select all zones.
	All ZoneFilterMode = "all"
	// Include is to select only a specific subset of zones.
	Include ZoneFilterMode = "include"
	// Exclude is to select all zones except a specific subset of zones.
	Exclude ZoneFilterMode = "exclude"
)

// LocationSelector represents a way to select zones from a given location.
type LocationSelector struct {
	// The unique name of the location.
	Name string
	// Define how the location zones will be selected.
	Mode ZoneFilterMode
	// List of zones which to be filtered in or out of the location zones, depending on the Mode.
	Zones []string
}

// Rule is an IDS rule, e.g. Snort or Suricata.
type Rule struct {
	// The unique rule ID.
	ID int64 `mutable:"false"`
	// The rule itself.
	Body string `mutable:"true"`
	// Select in which organization and zone the rule is enabled, e.g. "google:dmz".
	LocZones []string `mutable:"true"`
	// Last modified time of the message. Applied by the Store.
	LastModified string `mutable:"true"`
}

// SensorRequestType represents the type of sensor request message.
type SensorRequestType string

// Sensor request types as described in the sensor proto.
const (
	DeployRules SensorRequestType = "DeployRules"
	ReloadRules SensorRequestType = "ReloadRules"
)

// SensorRequest contains the details and state of a sensor request message.
type SensorRequest struct {
	// The request message ID.
	ID string `mutable:"false"`
	// The creation time of the message.
	Time string `mutable:"false"`
	// Fleetspeak client ID (Hex-encoded bytes).
	ClientID string `mutable:"false"`
	// Type of message.
	Type SensorRequestType `mutable:"false"`
	// Status of the request.
	Status string `mutable:"true"`
	// Last modified time of the message. Applied by the Store.
	LastModified string `mutable:"true"`
}

// SensorMessageType represents the type of message issued from a sensor.
type SensorMessageType string

const (
	// Response represents a sensor response to a sensor request.
	Response SensorMessageType = "Response"
	// Alert represents a sensor alert.
	Alert SensorMessageType = "Alert"
	// Heartbeat represents a sensor heartbeat.
	Heartbeat SensorMessageType = "Heartbeat"
)

// SensorMessage contains the details and state of a sensor message.
type SensorMessage struct {
	// The message ID.
	ID string `mutable:"false"`
	// The creation time of the message.
	Time string `mutable:"false"`
	// Fleetspeak client ID (Hex-encoded bytes).
	ClientID string `mutable:"false"`
	// Type of message.
	Type SensorMessageType `mutable:"false"`
	// Host information of sender.
	Host string `mutable:"false"`
	// Status of the request.
	Status string `mutable:"false"`
}
