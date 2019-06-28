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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"

	log "github.com/golang/glog"
	spb "github.com/google/emitto/source/sensor/proto"
	pb "github.com/google/emitto/source/server/proto"
)

// ProtoToLocation converts a proto Location to an internal Location.
func ProtoToLocation(l *pb.Location) *Location {
	var zones []string
	for _, z := range l.Zones {
		zones = append(zones, z)
	}
	return &Location{
		Name:  l.Name,
		Zones: zones,
	}
}

// LocationToProto converts an internal Location to proto Location.
func LocationToProto(l *Location) *pb.Location {
	var zones []string
	for _, z := range l.Zones {
		zones = append(zones, z)
	}
	return &pb.Location{
		Name:  l.Name,
		Zones: zones,
	}
}

// ProtoToRule converts a proto Rule to an internal Rule.
func ProtoToRule(r *pb.Rule) *Rule {
	var zones []string
	for _, z := range r.LocationZones {
		zones = append(zones, z)
	}
	return &Rule{
		ID:       r.Id,
		Body:     r.Body,
		LocZones: zones,
	}
}

// RuleToProto converts an internal Rule to a proto Rule.
func RuleToProto(r *Rule) *pb.Rule {
	var zones []string
	for _, z := range r.LocZones {
		zones = append(zones, z)
	}
	return &pb.Rule{
		Id:            r.ID,
		Body:          r.Body,
		LocationZones: zones,
	}
}

// MakeRuleFile builds a rule file given Rule objects.
func MakeRuleFile(rules []*Rule) []byte {
	var buf bytes.Buffer
	for _, r := range rules {
		buf.WriteString(fmt.Sprintf("%s\n", r.Body))
	}
	return buf.Bytes()
}

// MutationsMapping returns a map of fields and their mutability for Rule, Location,
// and SensorMessage objects.
//
// Fields are in the form "field_name" where "struct.FieldName" = "field_name".
// obj must not be a pointer.
func MutationsMapping(obj interface{}) (map[string]bool, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		return nil, errors.New("object must not be a pointer")
	}
	m := make(map[string]bool)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("mutable")
		if tag == "" {
			return nil, errors.New(`all Rule fields must contain a "mutable" tag`)
		}
		mutable, err := strconv.ParseBool(tag)
		if err != nil {
			return nil, err
		}
		m[strings.ToLower(strings.Join(camelcase.Split(v.Type().Field(i).Name), "_"))] = mutable
	}
	return m, nil
}

// ProtoToSensorMessage converts a proto sensor message to an internal SensorMessage.
func ProtoToSensorMessage(m *spb.SensorMessage) *SensorMessage {
	msg := &SensorMessage{
		ID: m.GetId(),
	}
	switch t := m.Type.(type) {
	case *spb.SensorMessage_Alert:
		msg.Time = time.Unix(m.GetAlert().GetTime().GetSeconds(), 0).Format(time.RFC1123Z)
		msg.Host = m.GetAlert().GetHost().String()
	case *spb.SensorMessage_Heartbeat:
		msg.Time = time.Unix(m.GetHeartbeat().GetTime().GetSeconds(), 0).Format(time.RFC1123Z)
		msg.Host = m.GetHeartbeat().GetHost().String()
	default:
		log.Errorf("Unknown sensor message type (%T)", t)
	}
	return msg
}

// ProtoToSensorRequest converts a proto SensorMessage to an internal SensorRequest.
func ProtoToSensorRequest(m *spb.SensorMessage) *SensorRequest {
	return &SensorRequest{
		ID:     m.GetResponse().GetId(),
		Status: m.GetResponse().GetStatus().String(),
	}
}
