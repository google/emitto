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
	"log"
	"reflect"
	"strings"
	"testing"
	"context"
	"time"

	"github.com/google/emitto/source/resources"
	"github.com/google/go-cmp/cmp"
)

// suite contains store tests and the underlying Store implementation.
type suite struct {
	builder func() (Store, error)
}

// RunTestSuite runs all generic store tests.
//
// The tests use the provided builder to instantiate a Store.
// The builder is expected to always return a valid Store.
func RunTestSuite(t *testing.T, builder func() (Store, error)) {
	s := &suite{builder}
	s.Run(t)
}

// Run runs Test* methods of the suite as subtests.
func (s *suite) Run(t *testing.T) {
	sv := reflect.ValueOf(s)
	st := reflect.TypeOf(s)
	for i := 0; i < sv.NumMethod(); i++ {
		n := st.Method(i).Name
		if strings.HasPrefix(n, "Test") {
			mv := sv.MethodByName(n)
			mt := mv.Type()
			if mt.NumIn() != 1 || !reflect.TypeOf(t).AssignableTo(mt.In(0)) {
				log.Fatalf("Method %q of the test suite must have 1 argument of type *testing.T", n)
			}
			if mt.NumOut() != 0 {
				log.Fatalf("Method %q of the test suite must have no return value", n)
			}
			m := mv.Interface().(func(t *testing.T))
			t.Run(n, m)
		}
	}
}

var (
	location1 = &resources.Location{
		Name:  "test",
		Zones: []string{"unstable", "canary", "prod"},
	}
	location2 = &resources.Location{
		Name:  "zzz_test",
		Zones: []string{"canary", "prod"},
	}

	rule1 = &resources.Rule{
		ID:   1111,
		Body: `sid:111 foo:bar`,
	}
	rule2 = &resources.Rule{
		ID:   2222,
		Body: `sid:222 foo:bar`,
	}
	rule3 = &resources.Rule{
		ID:   3333,
		Body: `sid:333 foo:bar`,
	}

	sensorRequest1 = &resources.SensorRequest{
		ID:       "req1",
		ClientID: "dest1",
		Type:     resources.DeployRules,
		Status:   "OK",
	}

	sensorMessage1 = &resources.SensorMessage{
		ID:       "req1",
		ClientID: "dest1",
		Type:     resources.Alert,
		Status:   "ERROR",
	}
)

func (s *suite) TestAddLocation(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddLocation(ctx, location1); err != nil {
		t.Error(err)
	}
	if err := st.AddLocation(ctx, location1); err == nil {
		t.Error("location.add: adding a duplicate location should have raised an error")
	}
}

func (s *suite) TestModifyLocation(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	TimeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}

	if err := st.AddLocation(ctx, location1); err != nil {
		t.Error(err)
	}
	cp := *location1
	cp.Zones = append(cp.Zones, "another_zone")
	if err := st.ModifyLocation(ctx, &cp); err != nil {
		t.Error(err)
	}
	got, err := st.GetLocation(ctx, cp.Name)
	if err != nil {
		t.Error(err)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	if diff := cmp.Diff(&cp, got); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func (s *suite) TestDeleteLocation(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddLocation(ctx, location1); err != nil {
		t.Error(err)
	}
	if err := st.DeleteLocation(ctx, location1.Name); err != nil {
		t.Error(err)
	}
	if _, err := st.GetLocation(ctx, location1.Name); err == nil {
		t.Error("GetLocation on a deleted location should have failed")
	}
}

func (s *suite) TestActionsOnNonExistingLocation(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.DeleteLocation(ctx, "does not exist"); err == nil {
		t.Error("DeleteLocation on a non-existing location should have failed")
	}
	if err := st.ModifyLocation(ctx, location1); err == nil {
		t.Error("ModifyLocation on a non-existing location should have failed")
	}
}

func (s *suite) TestListLocations(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	TimeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}

	if err := st.AddLocation(ctx, location1); err != nil {
		t.Error(err)
	}
	if err := st.AddLocation(ctx, location2); err != nil {
		t.Error(err)
	}
	got, err := st.ListLocations(ctx)
	if err != nil {
		t.Error(err)
	}

	location1.LastModified = TimeNow().Format(time.RFC1123Z)
	location2.LastModified = TimeNow().Format(time.RFC1123Z)
	want := []*resources.Location{location1, location2}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func (s *suite) TestAddRule(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddRule(ctx, rule1); err != nil {
		t.Error(err)
	}
	if err := st.AddRule(ctx, rule1); err == nil {
		t.Error("adding a duplicate rule should have raised an error")
	}
}

func (s *suite) TestDeleteRule(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddRule(ctx, rule1); err != nil {
		t.Error(err)
	}
	if err := st.DeleteRule(ctx, rule1.ID); err != nil {
		t.Error(err)
	}
	if _, err := st.ListRules(ctx, []int64{rule1.ID}); err == nil {
		t.Error("ListRules on a deleted rule should have failed")
	}
}

func (s *suite) TestActionsOnNonExistingRule(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.DeleteRule(ctx, -1); err == nil {
		t.Error("DeleteRule on a non-existing rule should have failed")
	}
	if err := st.ModifyRule(ctx, rule1); err == nil {
		t.Error("ModifyRule on a non-existing rule should have failed")
	}
}

func (s *suite) TestModifyRule(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	TimeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}

	if err := st.AddRule(ctx, rule1); err != nil {
		t.Error(err)
	}
	cp := *rule1
	cp.Body += " some_suffix:111"
	if err := st.ModifyRule(ctx, &cp); err != nil {
		t.Error(err)
	}
	got, err := st.ListRules(ctx, []int64{cp.ID})
	if err != nil {
		t.Error(err)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	if diff := cmp.Diff([]*resources.Rule{&cp}, got); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func (s *suite) TestListAllRules(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddRule(ctx, rule1); err != nil {
		t.Error(err)
	}
	if err := st.AddRule(ctx, rule2); err != nil {
		t.Error(err)
	}
	if err := st.AddRule(ctx, rule3); err != nil {
		t.Error(err)
	}
	got, err := st.ListRules(ctx, nil)
	if err != nil {
		t.Error(err)
	}
	if l := len(got); l != 3 {
		t.Errorf("expected 3 rules, got %d", l)
	}
}

func (s *suite) TestAddSensorRequest(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddSensorRequest(ctx, sensorRequest1); err != nil {
		t.Error(err)
	}
	if err := st.AddSensorRequest(ctx, sensorRequest1); err == nil {
		t.Error("adding a duplicate sensor request should have raised an error")
	}
}

func (s *suite) TestModifySensorRequest(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	TimeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}

	if err := st.AddSensorRequest(ctx, sensorRequest1); err != nil {
		t.Error(err)
	}
	cp := *sensorRequest1
	cp.Status = "NEW"
	if err := st.ModifySensorRequest(ctx, &cp); err != nil {
		t.Error(err)
	}
	got, err := st.GetSensorRequest(ctx, cp.ID)
	if err != nil {
		t.Error(err)
	}
	cp.LastModified = TimeNow().Format(time.RFC1123Z)
	if diff := cmp.Diff(&cp, got); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func (s *suite) TestDeleteSensorRequest(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddSensorRequest(ctx, sensorRequest1); err != nil {
		t.Error(err)
	}
	if err := st.DeleteSensorRequest(ctx, sensorRequest1.ID); err != nil {
		t.Error(err)
	}
	if _, err := st.GetSensorRequest(ctx, sensorRequest1.ID); err == nil {
		t.Error("GetSensorRequest on a deleted sensor request should have failed")
	}
}

func (s *suite) TestAddSensorMessage(t *testing.T) {
	st, err := s.builder()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if err := st.AddSensorMessage(ctx, sensorMessage1); err != nil {
		t.Error(err)
	}
	if err := st.AddSensorMessage(ctx, sensorMessage1); err == nil {
		t.Error("adding a duplicate sensor message should have raised an error")
	}
}
