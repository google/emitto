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

// Package store contains functionality to store objects.
package store

import (
	"context"

	"github.com/google/emitto/source/resources"
)

// Store represents object storage.
type Store interface {
	// AddLocation adds a new Location.
	AddLocation(ctx context.Context, l *resources.Location) error
	// ModifyLocation modifies an existing Location.
	ModifyLocation(ctx context.Context, l *resources.Location) error
	// DeleteLocation removes an existing Loction by name.
	DeleteLocation(ctx context.Context, name string) error
	// GetLocation retrieves a Location by name.
	GetLocation(ctx context.Context, name string) (*resources.Location, error)
	// ListLocations lists all stored Locations.
	ListLocations(ctx context.Context) ([]*resources.Location, error)

	// AddRule adds a new Rule.
	AddRule(ctx context.Context, r *resources.Rule) error
	// ModifyRule modifies an existing Rule.
	ModifyRule(ctx context.Context, r *resources.Rule) error
	// DeleteRule removes an existing Rule by ID.
	DeleteRule(ctx context.Context, id int64) error
	// ListRules lists stored Rules by ID.
	ListRules(ctx context.Context, ids []int64) ([]*resources.Rule, error)

	// AddSensorRequest adds a new SensorRequest.
	AddSensorRequest(ctx context.Context, r *resources.SensorRequest) error
	// ModifySensorRequest updates an existing SensorRequest.
	ModifySensorRequest(ctx context.Context, r *resources.SensorRequest) error
	// DeleteSensorRequest removes an existing SensorRequest by request ID.
	DeleteSensorRequest(ctx context.Context, id string) error
	// GetSensorRequest retrieves a SensorRequest by ID.
	GetSensorRequest(ctx context.Context, id string) (*resources.SensorRequest, error)

	// AddSensorMessage adds a new SensorMessage.
	AddSensorMessage(ctx context.Context, r *resources.SensorMessage) error
}
