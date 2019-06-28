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
	"time"
	"context"

	"github.com/google/emitto/source/resources"
	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

const (
  datastoreAddr     = "dns:///datastore.googleapis.com:443"
	locationKind      = "Location"
	ruleKind          = "Rule"
	sensorRequestKind = "SensorRequest"
	sensorMessageKind = "SensorMessage"
)

// DataStore represents a Google Cloud Datastore implementation of a Store.
type DataStore struct {
	client *datastore.Client
}

// NewDataStore returns a new DataStore.
func NewDataStore(client *datastore.Client) *DataStore {
	return &DataStore{client}
}

// NewGCDClient initializes a new Google Cloud Datastore Client.
// Follow these instructions to set up application credentials:
// https://cloud.google.com/docs/authentication/production#obtaining_and_providing_service_account_credentials_manually.
func NewGCDClient(ctx context.Context, projectID, credFile string) (*datastore.Client, error) {
	c, err := datastore.NewClient(ctx, projectID, option.WithEndpoint(datastoreAddr), option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, fmt.Errorf("GCD client creation failed: %v", err)
	}
	return c, nil
}

// Close Datastore client connection.
func (s *DataStore) Close() error {
	return s.client.Close()
}

func locationKey(name string) *datastore.Key {
	return &datastore.Key{
		Kind: locationKind,
		Name: name,
	}
}

// locationExists returns true if there is a location with the given name.
func (s *DataStore) locationExists(ctx context.Context, name string) (bool, error) {
	query := datastore.NewQuery(locationKind).Filter("__key__ =", locationKey(name)).KeysOnly()
	c, err := s.client.Count(ctx, query)
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

// AddLocation adds the given location.
func (s *DataStore) AddLocation(ctx context.Context, l *resources.Location) error {
	ok, err := s.locationExists(ctx, l.Name)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("location %q already exists", l.Name)
	}
	l.LastModified = TimeNow().Format(time.RFC1123Z)
	_, err = s.client.Put(ctx, locationKey(l.Name), l)
	return err
}

// ModifyLocation modifies an existing location with the provided location.
func (s *DataStore) ModifyLocation(ctx context.Context, l *resources.Location) error {
	ok, err := s.locationExists(ctx, l.Name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("location %q does not exist", l.Name)
	}
	loc, err := s.GetLocation(ctx, l.Name)
	if err != nil {
		return fmt.Errorf("unable to get location %q: %v", l.Name, err)
	}
	if err := MutateLocation(l, loc); err != nil {
		return fmt.Errorf("unable to mutate location src=%+v dst=%+v: %v", l, loc, err)
	}
	loc.LastModified = TimeNow().Format(time.RFC1123Z)
	_, err = s.client.Put(ctx, locationKey(l.Name), loc)
	return err
}

// DeleteLocation deletes the given location.
func (s *DataStore) DeleteLocation(ctx context.Context, name string) error {
	ok, err := s.locationExists(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("location %q does not exist", name)
	}
	return s.client.Delete(ctx, locationKey(name))
}

// GetLocation gets the location with the given name.
func (s *DataStore) GetLocation(ctx context.Context, name string) (*resources.Location, error) {
	query := datastore.NewQuery(locationKind).Filter("__key__ =", locationKey(name))
	l := new(resources.Location)
	if _, err := s.client.Run(ctx, query).Next(l); err != nil {
		return nil, err
	}
	return l, nil
}

// ListLocations list all locations, ordered by name.
func (s *DataStore) ListLocations(ctx context.Context) ([]*resources.Location, error) {
	var all []*resources.Location
	query := datastore.NewQuery(locationKind).Order("Name")
	if _, err := s.client.GetAll(ctx, query, &all); err != nil {
		return nil, err
	}
	return all, nil
}

func ruleKey(ruleID int64) *datastore.Key {
	return &datastore.Key{
		Kind: ruleKind,
		ID:   ruleID,
	}
}

// ruleExists returns trues if there is a rule with the given rule ID.
func (s *DataStore) ruleExists(ctx context.Context, id int64) (bool, error) {
	query := datastore.NewQuery(ruleKind).Filter("__key__ =", ruleKey(id)).KeysOnly()
	c, err := s.client.Count(ctx, query)
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

// AddRule adds the given rule.
func (s *DataStore) AddRule(ctx context.Context, r *resources.Rule) error {
	switch ok, err := s.ruleExists(ctx, r.ID); {
	case err != nil:
		return err
	case ok:
		return fmt.Errorf("rule %d already exists", r.ID)
	default:
		r.LastModified = TimeNow().Format(time.RFC1123Z)
		_, err = s.client.Put(ctx, ruleKey(r.ID), r)
		return err
	}
}

// ModifyRule modifies an existing rule with the provided rule.
func (s *DataStore) ModifyRule(ctx context.Context, r *resources.Rule) error {
	ok, err := s.ruleExists(ctx, r.ID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("rule %d does not exist", r.ID)
	}
	rules, err := s.ListRules(ctx, []int64{r.ID})
	if err != nil {
		return fmt.Errorf("unable to get rule %d: %v", r.ID, err)
	}
	if l := len(rules); l > 1 {
		return fmt.Errorf("expected 1 rule (ID=%d), got %d", r.ID, l)
	}
	rule := rules[0]
	if err := MutateRule(r, rule); err != nil {
		return fmt.Errorf("unable to mutate rule src=%+v dst=%+v: %v", r, rule, err)
	}
	rule.LastModified = TimeNow().Format(time.RFC1123Z)
	_, err = s.client.Put(ctx, ruleKey(r.ID), rule)
	return err
}

// DeleteRule deletes the given rule.
func (s *DataStore) DeleteRule(ctx context.Context, id int64) error {
	switch ok, err := s.ruleExists(ctx, id); {
	case err != nil:
		return err
	case !ok:
		return fmt.Errorf("rule %d does not exist", id)
	default:
		return s.client.Delete(ctx, ruleKey(id))
	}
}

// ListRules lists the rules with the given rule IDs, following the same order.
// If `ids` is nil, lists all rules, ordered by ID.
func (s *DataStore) ListRules(ctx context.Context, ids []int64) ([]*resources.Rule, error) {
	var (
		all []*resources.Rule
		err error
	)
	if l := len(ids); l > 0 {
		all = make([]*resources.Rule, l)
		keys := make([]*datastore.Key, l)
		for j, id := range ids {
			keys[j] = ruleKey(id)
		}
		err = s.client.GetMulti(ctx, keys, all)
	} else {
		query := datastore.NewQuery(ruleKind).Order("ID")
		if _, err := s.client.GetAll(ctx, query, &all); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return all, nil
}

func sensorRequestKey(id string) *datastore.Key {
	return &datastore.Key{
		Kind: sensorRequestKind,
		Name: id,
	}
}

// sensorRequestExists returns true if there is a sensor request with the given ID.
func (s *DataStore) sensorRequestExists(ctx context.Context, id string) (bool, error) {
	query := datastore.NewQuery(sensorRequestKind).Filter("__key__ =", sensorRequestKey(id)).KeysOnly()
	c, err := s.client.Count(ctx, query)
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

// AddSensorRequest adds the given sensor request.
func (s *DataStore) AddSensorRequest(ctx context.Context, r *resources.SensorRequest) error {
	switch ok, err := s.sensorRequestExists(ctx, r.ID); {
	case err != nil:
		return err
	case ok:
		return fmt.Errorf("sensor request %q already exists", r.ID)
	default:
		r.LastModified = TimeNow().Format(time.RFC1123Z)
		_, err = s.client.Put(ctx, sensorRequestKey(r.ID), r)
		return err
	}
}

// ModifySensorRequest modifies an existing sensor request with the provided sensor request.
func (s *DataStore) ModifySensorRequest(ctx context.Context, r *resources.SensorRequest) error {
	ok, err := s.sensorRequestExists(ctx, r.ID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("sensor request %q does not exist", r.ID)
	}
	req, err := s.GetSensorRequest(ctx, r.ID)
	if err != nil {
		return fmt.Errorf("unable to get sensor request %q: %v", r.ID, err)
	}
	if err := MutateSensorRequest(r, req); err != nil {
		return fmt.Errorf("unable to mutate sensor request src=%+v dst=%+v: %v", r, req, err)
	}
	req.LastModified = TimeNow().Format(time.RFC1123Z)
	_, err = s.client.Put(ctx, sensorRequestKey(r.ID), req)
	return err
}

// DeleteSensorRequest removes the given sensor request.
func (s *DataStore) DeleteSensorRequest(ctx context.Context, id string) error {
	switch ok, err := s.sensorRequestExists(ctx, id); {
	case err != nil:
		return err
	case !ok:
		return fmt.Errorf("sensor request %q does not exist", id)
	default:
		return s.client.Delete(ctx, sensorRequestKey(id))
	}
}

// GetSensorRequest gets the sensor request with the given ID.
func (s *DataStore) GetSensorRequest(ctx context.Context, id string) (*resources.SensorRequest, error) {
	query := datastore.NewQuery(sensorRequestKind).Filter("__key__ =", sensorRequestKey(id))
	l := new(resources.SensorRequest)
	if _, err := s.client.Run(ctx, query).Next(l); err != nil {
		return nil, err
	}
	return l, nil
}

func sensorMessageKey(id string) *datastore.Key {
	return &datastore.Key{
		Kind: sensorMessageKind,
		Name: id,
	}
}

// sensorMessageExists returns true if there is a sensor request with the given ID.
func (s *DataStore) sensorMessageExists(ctx context.Context, id string) (bool, error) {
	query := datastore.NewQuery(sensorMessageKind).Filter("__key__ =", sensorMessageKey(id)).KeysOnly()
	c, err := s.client.Count(ctx, query)
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

// AddSensorMessage adds the given sensor message.
func (s *DataStore) AddSensorMessage(ctx context.Context, m *resources.SensorMessage) error {
	switch ok, err := s.sensorMessageExists(ctx, m.ID); {
	case err != nil:
		return err
	case ok:
		return fmt.Errorf("sensor request %q already exists", m.ID)
	default:
		_, err = s.client.Put(ctx, sensorMessageKey(m.ID), m)
		return err
	}
}
