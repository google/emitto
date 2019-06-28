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

package store

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"context"

	"github.com/google/emitto/source/resources"
)

// MemoryStore represents a memory Store implementation.
type MemoryStore struct {
	m              sync.Mutex
	locations      map[string]resources.Location
	rules          map[int64]resources.Rule
	sensorRequests map[string]resources.SensorRequest
	sensorMessages map[string]resources.SensorMessage
}

// NewMemoryStore returns a MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		locations:      make(map[string]resources.Location),
		rules:          make(map[int64]resources.Rule),
		sensorRequests: make(map[string]resources.SensorRequest),
		sensorMessages: make(map[string]resources.SensorMessage),
	}
}

// AddLocation adds a Location to the store.
func (s *MemoryStore) AddLocation(ctx context.Context, l *resources.Location) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *l
	if _, ok := s.locations[cp.Name]; ok {
		return fmt.Errorf("location %q already exists", cp.Name)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	s.locations[cp.Name] = cp
	return nil
}

// ModifyLocation modifies an existing location with the provided location.
func (s *MemoryStore) ModifyLocation(ctx context.Context, l *resources.Location) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *l
	loc, ok := s.locations[cp.Name]
	if !ok {
		return fmt.Errorf("location %q does not exist", cp.Name)
	}
	if err := MutateLocation(&cp, &loc); err != nil {
		return fmt.Errorf("unable to mutate location src=%+v dst=%+v: %v", cp, loc, err)
	}
	loc.LastModified = TimeNow().Format(time.RFC1123Z)
	s.locations[cp.Name] = loc
	return nil
}

// DeleteLocation deletes an existing Location.
func (s *MemoryStore) DeleteLocation(ctx context.Context, name string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.locations[name]; !ok {
		return fmt.Errorf("location %q does not exist", name)
	}
	delete(s.locations, name)
	return nil
}

// GetLocation returns the Location with the given name.
func (s *MemoryStore) GetLocation(ctx context.Context, name string) (*resources.Location, error) {
	s.m.Lock()
	defer s.m.Unlock()

	r, ok := s.locations[name]
	if !ok {
		return nil, fmt.Errorf("location %q does not exist", name)
	}
	return &r, nil
}

// ListLocations returns all the locations, sorted by name.
func (s *MemoryStore) ListLocations(ctx context.Context) ([]*resources.Location, error) {
	s.m.Lock()
	defer s.m.Unlock()

	locations := make([]*resources.Location, 0, len(s.locations))
	for l := range s.locations {
		loc := s.locations[l]
		locations = append(locations, &loc)
	}
	sort.Slice(locations, func(i, j int) bool {
		return locations[i].Name < locations[j].Name
	})
	return locations, nil
}

// AddRule adds a Rule to the store.
func (s *MemoryStore) AddRule(ctx context.Context, r *resources.Rule) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *r
	if _, ok := s.rules[cp.ID]; ok {
		return fmt.Errorf("rule %d already exists", cp.ID)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	s.rules[cp.ID] = cp
	return nil
}

// ModifyRule modifies an existing rule with the provided rule.
func (s *MemoryStore) ModifyRule(ctx context.Context, r *resources.Rule) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *r
	rule, ok := s.rules[cp.ID]
	if !ok {
		return fmt.Errorf("rule %d does not exist", cp.ID)
	}
	if err := MutateRule(&cp, &rule); err != nil {
		return fmt.Errorf("unable to mutate rule src=%+v dst=%+v: %v", cp, rule, err)
	}
	rule.LastModified = TimeNow().Format(time.RFC1123Z)
	s.rules[cp.ID] = rule
	return nil
}

// DeleteRule deletes an existing Rule.
func (s *MemoryStore) DeleteRule(ctx context.Context, id int64) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.rules[id]; !ok {
		return fmt.Errorf("rule %d does not exist", id)
	}
	delete(s.rules, id)
	return nil
}

// ListRules returns Rules from a list of rule IDs. All rules are returned if ids is nil.
func (s *MemoryStore) ListRules(ctx context.Context, ids []int64) ([]*resources.Rule, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(ids) == 0 {
		var rules []*resources.Rule
		for r := range s.rules {
			rule := s.rules[r]
			rules = append(rules, &rule)
		}
		return rules, nil
	}

	rules := make([]*resources.Rule, 0, len(ids))
	for _, id := range ids {
		r, ok := s.rules[id]
		if !ok {
			return nil, fmt.Errorf("rule %d does not exist", id)
		}
		rules = append(rules, &r)
	}
	return rules, nil
}

// AddSensorRequest adds a sensor request.
func (s *MemoryStore) AddSensorRequest(ctx context.Context, r *resources.SensorRequest) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *r
	if _, ok := s.sensorRequests[cp.ID]; ok {
		return fmt.Errorf("sensor request %q already exists", cp.ID)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	s.sensorRequests[cp.ID] = cp
	return nil
}

// ModifySensorRequest modifies an existing sensor request with the provided sensor request.
func (s *MemoryStore) ModifySensorRequest(ctx context.Context, r *resources.SensorRequest) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *r
	req, ok := s.sensorRequests[cp.ID]
	if !ok {
		return fmt.Errorf("sensor request %q does not exist", cp.ID)
	}
	if err := MutateSensorRequest(&cp, &req); err != nil {
		return fmt.Errorf("unable to mutate sensor request src=%+v dst=%+v: %v", cp, req, err)
	}
	req.LastModified = TimeNow().Format(time.RFC1123Z)
	s.sensorRequests[cp.ID] = req
	return nil
}

// DeleteSensorRequest deletes an existing sensor request.
func (s *MemoryStore) DeleteSensorRequest(ctx context.Context, id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.sensorRequests[id]; !ok {
		return fmt.Errorf("sensor request %q does not exist", id)
	}
	delete(s.sensorRequests, id)
	return nil
}

// GetSensorRequest returns the sensor request with the given ID.
func (s *MemoryStore) GetSensorRequest(ctx context.Context, id string) (*resources.SensorRequest, error) {
	s.m.Lock()
	defer s.m.Unlock()

	r, ok := s.sensorRequests[id]
	if !ok {
		return nil, fmt.Errorf("sensor request %q does not exist", id)
	}
	return &r, nil
}

// AddSensorMessage adds a sensor message.
func (s *MemoryStore) AddSensorMessage(ctx context.Context, m *resources.SensorMessage) error {
	s.m.Lock()
	defer s.m.Unlock()

	cp := *m
	if _, ok := s.sensorMessages[cp.ID]; ok {
		return fmt.Errorf("sensor message %q already exists", cp.ID)
	}
	s.sensorMessages[cp.ID] = cp
	return nil
}
