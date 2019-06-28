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
	"testing"

	"github.com/google/emitto/source/resources"
	"github.com/google/go-cmp/cmp"
)

func TestMutateRule(t *testing.T) {
	for _, tt := range []struct {
		desc           string
		src, dst, want *resources.Rule
		overrideObj    bool
	}{
		{
			desc: "Body is mutated",
			src: &resources.Rule{
				Body: "new_body",
			},
			dst: &resources.Rule{
				ID:   2222,
				Body: "old_body",
			},
			want: &resources.Rule{
				ID:   2222,
				Body: "new_body",
			},
		},
		{
			desc: "ID is not mutated",
			src: &resources.Rule{
				ID: 1111,
			},
			dst: &resources.Rule{
				ID: 2222,
			},
			want: &resources.Rule{
				ID: 2222,
			},
		},
		{
			desc: "Body and LocZones are mutated",
			src: &resources.Rule{
				ID:       2222,
				Body:     "new_body",
				LocZones: []string{"new_zone_1", "new_zone_2"},
			},
			dst: &resources.Rule{
				ID:       2222,
				Body:     "old_body",
				LocZones: []string{"old_zone_1", "old_zone_2"},
			},
			want: &resources.Rule{
				ID:       2222,
				Body:     "new_body",
				LocZones: []string{"new_zone_1", "new_zone_2"},
			},
		},
	} {
		if err := MutateRule(tt.src, tt.dst); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if diff := cmp.Diff(tt.want, tt.dst); diff != "" {
			t.Errorf("%s: expectation mismatch (-want +got):\n%s", tt.desc, diff)
		}
	}
}

func TestMutateLocation(t *testing.T) {
	for _, tt := range []struct {
		desc           string
		src, dst, want *resources.Location
	}{
		{
			desc: "Zones are mutated",
			src: &resources.Location{
				Zones: []string{"new_zone_1", "new_zone_2"},
			},
			dst: &resources.Location{
				Zones: []string{"old_zone_1", "old_zone_2"},
			},
			want: &resources.Location{
				Zones: []string{"new_zone_1", "new_zone_2"},
			},
		},
		{
			desc: "Name is not mutated",
			src: &resources.Location{
				Name: "new_name",
			},
			dst: &resources.Location{
				Name: "old_name",
			},
			want: &resources.Location{
				Name: "old_name",
			},
		},
	} {
		if err := MutateLocation(tt.src, tt.dst); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if diff := cmp.Diff(tt.want, tt.dst); diff != "" {
			t.Errorf("%s: expectation mismatch (-want +got):\n%s", tt.desc, diff)
		}
	}
}

func TestMutateSensorRequest(t *testing.T) {
	for _, tt := range []struct {
		desc           string
		src, dst, want *resources.SensorRequest
	}{
		{
			desc: "Status is mutated",
			src: &resources.SensorRequest{
				Status: "new",
			},
			dst: &resources.SensorRequest{
				Status: "old",
			},
			want: &resources.SensorRequest{
				Status: "new",
			},
		},
		{
			desc: "preserve non mutable fields",
			src: &resources.SensorRequest{
				ID:       "new_id",
				Time:     "new_time",
				ClientID: "new_client_id",
				Type:     "new_type",
			},
			dst: &resources.SensorRequest{
				ID:       "old_id",
				Time:     "old_time",
				ClientID: "old_client_id",
				Type:     "old_type",
			},
			want: &resources.SensorRequest{
				ID:       "old_id",
				Time:     "old_time",
				ClientID: "old_client_id",
				Type:     "old_type",
			},
		},
	} {
		if err := MutateSensorRequest(tt.src, tt.dst); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if diff := cmp.Diff(tt.want, tt.dst); diff != "" {
			t.Errorf("%s: expectation mismatch (-want +got):\n%s", tt.desc, diff)
		}
	}
}
