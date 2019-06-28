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

package service

import (
	"testing"

	"github.com/google/emitto/source/resources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	spb "github.com/google/emitto/source/server/proto"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
	fsspb "github.com/google/fleetspeak/fleetspeak/src/server/proto/fleetspeak_server"
	mpb "google.golang.org/genproto/protobuf/field_mask"
)

var (
	testClients = []*fsspb.Client{
		{
			ClientId:        []byte("client_a"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-a"}, {Label: "alphabet-location-zone-dmz"}},
			LastContactTime: &tspb.Timestamp{Seconds: 1111111111},
		},
		{
			ClientId:        []byte("client_b"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-a"}, {Label: "alphabet-location-zone-dmz"}},
			LastContactTime: &tspb.Timestamp{Seconds: 2222222222},
		},
		{
			ClientId:        []byte("client_c"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-a"}, {Label: "alphabet-location-zone-corp"}},
			LastContactTime: &tspb.Timestamp{Seconds: 3333333333},
		},
		{
			ClientId:        []byte("client_d"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-b"}, {Label: "alphabet-location-zone-dmz"}},
			LastContactTime: &tspb.Timestamp{Seconds: 4444444444},
		},
		{
			ClientId:        []byte("client_e"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-c"}, {Label: "alphabet-location-zone-corp"}},
			LastContactTime: &tspb.Timestamp{Seconds: 5555555555},
		},
		{
			ClientId:        []byte("client_f"),
			Labels:          []*fspb.Label{{Label: "alphabet-location-name-a"}, {Label: "alphabet-location-zone-prod"}},
			LastContactTime: &tspb.Timestamp{Seconds: 6666666666},
		},
	}
	testRules = []*resources.Rule{
		{
			ID:       1111,
			Body:     "test",
			LocZones: []string{"a:dmz", "b:corp"},
		},
		{
			ID:       2222,
			Body:     "test",
			LocZones: []string{"a:corp", "b:dmz"},
		},
		{
			ID:       3333,
			Body:     "test",
			LocZones: []string{"c:prod"},
		},
		{
			ID:       4444,
			Body:     "test",
			LocZones: []string{"b:dmz", "b:corp"},
		},
		{
			ID:       5555,
			Body:     "test",
			LocZones: []string{"google:dmz"},
		},
	}
	testLocations = []*resources.Location{
		{
			Name:  "a",
			Zones: []string{"a"},
		},
		{
			Name:  "b",
			Zones: []string{"a", "b"},
		},
		{
			Name:  "c",
			Zones: []string{"a", "b", "c"},
		},
	}
)

func TestGetClientIDsByLocation(t *testing.T) {
	tests := []struct {
		desc string
		loc  *spb.Location
		want []string
	}{
		{
			desc: "1 zone",
			loc:  &spb.Location{Name: "a", Zones: []string{"dmz"}},
			want: []string{"client_a", "client_b"},
		},
		{
			desc: "2 zones",
			loc:  &spb.Location{Name: "a", Zones: []string{"dmz", "corp"}},
			want: []string{"client_a", "client_b", "client_c"},
		},
		{
			desc: "3 zones",
			loc:  &spb.Location{Name: "a", Zones: []string{"dmz", "corp", "prod"}},
			want: []string{"client_a", "client_b", "client_c", "client_f"},
		},
		{
			desc: "non-existent location",
			loc:  &spb.Location{Name: "d", Zones: []string{"corp"}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got []string
			for _, id := range getClientIDsByLocation(testClients, tt.loc) {
				got = append(got, string(id))
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("expectation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFilterRulesByLocation(t *testing.T) {
	tests := []struct {
		desc string
		loc  *spb.Location
		want []*resources.Rule
	}{
		{
			desc: "1 rule",
			loc:  &spb.Location{Name: "a", Zones: []string{"dmz"}},
			want: []*resources.Rule{{ID: 1111, Body: "test", LocZones: []string{"a:dmz", "b:corp"}}},
		},
		{
			desc: "3 testRules",
			loc:  &spb.Location{Name: "b", Zones: []string{"dmz", "corp"}},
			want: []*resources.Rule{
				{
					ID:       1111,
					Body:     "test",
					LocZones: []string{"a:dmz", "b:corp"},
				},
				{
					ID:       2222,
					Body:     "test",
					LocZones: []string{"a:corp", "b:dmz"},
				},
				{
					ID:       4444,
					Body:     "test",
					LocZones: []string{"b:dmz", "b:corp"},
				},
			},
		},
		{
			desc: "unknown loczone",
			loc:  &spb.Location{Name: "goog", Zones: []string{"dmz"}},
			want: []*resources.Rule{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := filterRulesByLocation(testRules, tt.loc)
			if diff := cmp.Diff(tt.want, got, cmpopts.SortSlices(func(a, b *resources.Rule) bool { return a.ID < b.ID })); diff != "" {
				t.Errorf("expectation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateUpdateMask(t *testing.T) {
	tests := []struct {
		desc      string
		obj       interface{}
		fieldMask *mpb.FieldMask
		wantErr   bool
	}{
		{
			desc: "valid Rule and mask paths",
			obj:  resources.Rule{},
			fieldMask: &mpb.FieldMask{
				Paths: []string{"body", "loc_zones"},
			},
		},
		{
			desc: "valid Rule but invalid mask path",
			obj:  resources.Rule{},
			fieldMask: &mpb.FieldMask{
				Paths: []string{"id"},
			},
			wantErr: true,
		},
		{
			desc: "invalid object",
			obj:  resources.LocationSelector{},
			fieldMask: &mpb.FieldMask{
				Paths: []string{"foo"},
			},
			wantErr: true,
		},
		{
			desc: "valid Location and valid mask paths",
			obj:  resources.Location{},
			fieldMask: &mpb.FieldMask{
				Paths: []string{"zones"},
			},
		},
		{
			desc: "valid Location but invalid mask path",
			obj:  resources.Rule{},
			fieldMask: &mpb.FieldMask{
				Paths: []string{"name"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if err := ValidateUpdateMask(tt.obj, tt.fieldMask); (err != nil) != tt.wantErr {
				t.Errorf("got err=%v, wantErr=%t", err, tt.wantErr)
			}
		})
	}
}
