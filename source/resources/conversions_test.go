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

package resources

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	tpb "github.com/golang/protobuf/ptypes/timestamp"
	spb "github.com/google/emitto/source/sensor/proto"
	pb "github.com/google/emitto/source/server/proto"
)

func TestProtoToLocation(t *testing.T) {
	p := &pb.Location{
		Name:  "test",
		Zones: []string{"dmz", "prod"},
	}
	want := &Location{
		Name:  "test",
		Zones: []string{"dmz", "prod"},
	}
	if diff := cmp.Diff(want, ProtoToLocation(p)); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestLocationToProto(t *testing.T) {
	l := &Location{
		Name:  "test",
		Zones: []string{"dmz", "prod"},
	}
	want := &pb.Location{
		Name:  "test",
		Zones: []string{"dmz", "prod"},
	}
	if diff := cmp.Diff(want, LocationToProto(l), cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestProtoToRule(t *testing.T) {
	p := &pb.Rule{
		Id:            1234567890,
		Body:          "test rule",
		LocationZones: []string{"test:dmz", "test:corp"},
	}
	want := &Rule{
		ID:       1234567890,
		Body:     "test rule",
		LocZones: []string{"test:dmz", "test:corp"},
	}
	if diff := cmp.Diff(want, ProtoToRule(p)); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestRuleToProto(t *testing.T) {
	r := &Rule{
		ID:       1234567890,
		Body:     "test rule",
		LocZones: []string{"test:dmz", "test:corp"},
	}
	want := &pb.Rule{
		Id:            1234567890,
		Body:          "test rule",
		LocationZones: []string{"test:dmz", "test:corp"},
	}
	if diff := cmp.Diff(want, RuleToProto(r), cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestMakeRuleFile(t *testing.T) {
	rules := []*Rule{
		{
			ID:       123,
			Body:     "test rule",
			LocZones: []string{"test:dmz", "test:corp"},
		},
		{
			ID:       1234,
			Body:     "test rule",
			LocZones: []string{"test:dmz", "test:corp"},
		},
		{
			ID:       12345,
			Body:     "test rule",
			LocZones: []string{"test:dmz", "test:corp"},
		},
	}
	want := []byte("test rule\ntest rule\ntest rule\n")
	if diff := cmp.Diff(string(want), string(MakeRuleFile(rules))); diff != "" {
		t.Errorf("expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestMutationsMapping(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		object  interface{}
		want    map[string]bool
		wantErr bool
	}{
		{
			desc: "valid mappings",
			object: struct {
				f1     string `mutable:"false"`
				f2, f3 string `mutable:"true"`
			}{},
			want: map[string]bool{"f_1": false, "f_2": true, "f_3": true},
		},
		{
			desc:    "missing tags",
			object:  struct{ f1 string }{},
			wantErr: true,
		},
		{
			desc:    "pointer object",
			object:  &struct{ f1 string }{},
			wantErr: true,
		},
	} {
		got, err := MutationsMapping(tt.object)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: got err=%v, wantErr=%v", tt.desc, err, tt.wantErr)
		}
		if err != nil {
			continue
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("%s: expectation mismatch (-want +got):\n%s", tt.desc, diff)
		}
	}
}

func TestProtoToSensorMessage(t *testing.T) {
	for _, tt := range []struct {
		desc string
		p    *spb.SensorMessage
		want *SensorMessage
	}{
		{
			desc: "heartbeat mismatch",
			p: &spb.SensorMessage{
				Id: "test_id",
				Type: &spb.SensorMessage_Heartbeat{
					Heartbeat: &spb.Heartbeat{
						Time: &tpb.Timestamp{Seconds: 123},
						Host: &spb.Host{Fqdn: "id1", Ip: "id2", Uuid: "id3", Org: "org", Zone: "zone"},
					},
				},
			},
			want: &SensorMessage{
				ID:   "test_id",
				Time: "Thu, 01 Jan 1970 00:02:03 +0000",
				Host: `fqdn:"id1" ip:"id2" uuid:"id3" org:"org" zone:"zone" `,
			},
		},
		{
			desc: "alert mistmatch",
			p: &spb.SensorMessage{
				Id: "test_id",
				Type: &spb.SensorMessage_Alert{
					Alert: &spb.SensorAlert{
						Time: &tpb.Timestamp{Seconds: 123},
						Host: &spb.Host{Fqdn: "id1", Ip: "id2", Uuid: "id3", Org: "org", Zone: "zone"},
					},
				},
			},
			want: &SensorMessage{
				ID:   "test_id",
				Time: "Thu, 01 Jan 1970 00:02:03 +0000",
				Host: `fqdn:"id1" ip:"id2" uuid:"id3" org:"org" zone:"zone" `,
			},
		},
	} {
		if diff := cmp.Diff(tt.want, ProtoToSensorMessage(tt.p)); diff != "" {
			t.Errorf("%s (-want +got):\n%s", tt.desc, diff)
		}
	}
}
