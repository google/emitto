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

// Package service provides the implementation of the Emitto service.
package service

import (
	"fmt"
	"time"
	"context"

	"github.com/google/emitto/source/filestore"
	"github.com/google/emitto/source/resources"
	"github.com/google/emitto/source/server/fleetspeak"
	"github.com/google/emitto/source/server/store"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/golang/glog"
	spb "github.com/google/emitto/source/sensor/proto"
	svpb "github.com/google/emitto/source/server/proto"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
	fsspb "github.com/google/fleetspeak/fleetspeak/src/server/proto/fleetspeak_server"
)

// FleetspeakAdminClient represents a Fleetspeak admin client.
type FleetspeakAdminClient interface {
	// Insert a SensorRequest message for delivery.
	InsertMessage(ctx context.Context, req *spb.SensorRequest, id []byte) error
	// List all clients.
	ListClients(ctx context.Context) ([]*fsspb.Client, error)
	// Close the RPC connection.
	Close() error
}

// Service contains handlers for Emitto server functions.
type Service struct {
	store      store.Store
	fileStore  filestore.FileStore
	fleetspeak FleetspeakAdminClient
}

// New returns a new emitto Service.
func New(store store.Store, filestore filestore.FileStore, fs FleetspeakAdminClient) *Service {
	return &Service{store, filestore, fs}
}

// DeployRules generates a rule file and deploys it to the sensors in the provided location.
func (s *Service) DeployRules(req *svpb.DeployRulesRequest, stream svpb.Emitto_DeployRulesServer) error {
	ctx := stream.Context()
	// Get clients.
	clients, err := s.fleetspeak.ListClients(ctx)
	if err != nil {
		return err
	}
	log.V(1).Infof("DeployRules() listed clients:\n%s", fleetspeak.ParseClients(clients))
	ids := getClientIDsByLocation(clients, req.GetLocation())
	if len(ids) == 0 {
		return status.Errorf(codes.FailedPrecondition, "no clients for location: %v", req.GetLocation())
	}
	// Create rule file.
	rules, err := s.store.ListRules(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list rules: %v", err)
	}

	if rules = filterRulesByLocation(rules, req.GetLocation()); len(rules) == 0 {
		return status.Errorf(codes.FailedPrecondition, "no rules found for %q", req.GetLocation())
	}
	path := ruleFilepath(req.GetLocation().GetName())
	if err := s.fileStore.AddRuleFile(ctx, path, resources.MakeRuleFile(rules)); err != nil {
		return err
	}

	for _, id := range ids {
		resp := &svpb.DeployRulesResponse{
			ClientId: fmt.Sprintf("%X", id),
			Status:   status.New(codes.OK, "OK").Proto(), // Default; will be reset for errors.
		}
		rid := uuid.New().String()
		// Log sensor request. This is a precondition to sending the request.
		m := &resources.SensorRequest{
			ID:       rid,
			Time:     timeNow().Format(time.RFC1123Z),
			ClientID: fmt.Sprintf("%X", id),
			Type:     resources.DeployRules,
		}
		if err := s.store.AddSensorRequest(ctx, m); err != nil {
			resp.Status = status.New(codes.FailedPrecondition, fmt.Sprintf("failed to add sensor message (%+v): %v", m, err)).Proto()
			if err := stream.Send(resp); err != nil {
				return err
			}
			continue
		}
		// Send sensor request to client.
		r := &spb.SensorRequest{
			Id:   rid,
			Time: &tspb.Timestamp{Seconds: time.Now().Unix()},
			Type: &spb.SensorRequest_DeployRules{&spb.DeployRules{RuleFile: path}},
		}
		if err := s.fleetspeak.InsertMessage(ctx, r, id); err != nil {
			// Clean up Store entry - do not fail hard on this.
			if err := s.store.DeleteSensorRequest(ctx, rid); err != nil {
				log.Errorf("Failed to remove sensor request (%+v): %v", r, err)
			}
			resp.Status = status.New(codes.Internal, fmt.Sprintf("failed to insert message: %v", err)).Proto()
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

// AddRule adds the provided Rule.
func (s *Service) AddRule(ctx context.Context, req *svpb.AddRuleRequest) (*emptypb.Empty, error) {
	if err := s.store.AddRule(ctx, resources.ProtoToRule(req.GetRule())); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to add rule: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ModifyRule modifies an existing Rule with the provided field mask.
func (s *Service) ModifyRule(ctx context.Context, req *svpb.ModifyRuleRequest) (*emptypb.Empty, error) {
	r := resources.ProtoToRule(req.GetRule())
	if err := ValidateUpdateMask(*r, req.GetFieldMask()); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument,
			"failed to validate modifications for rule (%+v) and mask (%+v): %v", r, req.GetFieldMask(), err)
	}
	if err := s.store.ModifyRule(ctx, r); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to modify rule (%+v): %v", r, err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteRule deletes an existing Rule by Rule ID.
func (s *Service) DeleteRule(ctx context.Context, req *svpb.DeleteRuleRequest) (*emptypb.Empty, error) {
	if err := s.store.DeleteRule(ctx, req.GetRuleId()); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete rule (id=%d): %v", req.GetRuleId(), err)
	}
	return &emptypb.Empty{}, nil
}

// ListRules returns Rules for the provided Rule IDs.
func (s *Service) ListRules(ctx context.Context, req *svpb.ListRulesRequest) (*svpb.ListRulesResponse, error) {
	rules, err := s.store.ListRules(ctx, req.GetRuleIds())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list rules (%+v): %v", req.GetRuleIds(), err)
	}
	resp := &svpb.ListRulesResponse{}
	for _, r := range rules {
		resp.Rules = append(resp.Rules, resources.RuleToProto(r))
	}
	return resp, nil
}

// AddLocation adds the provided Location.
func (s *Service) AddLocation(ctx context.Context, req *svpb.AddLocationRequest) (*emptypb.Empty, error) {
	if err := s.store.AddLocation(ctx, resources.ProtoToLocation(req.GetLocation())); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to add location (%+v): %v", req.GetLocation(), err)
	}
	return &emptypb.Empty{}, nil
}

// ModifyLocation updates the specified location.
func (s *Service) ModifyLocation(ctx context.Context, req *svpb.ModifyLocationRequest) (*emptypb.Empty, error) {
	l := resources.ProtoToLocation(req.GetLocation())
	if err := ValidateUpdateMask(*l, req.GetFieldMask()); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "failed to modify location (%+v): %v", l, err)
	}
	if err := s.store.ModifyLocation(ctx, l); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to modify location: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteLocation deletes an existing Location by Location Name.
func (s *Service) DeleteLocation(ctx context.Context, req *svpb.DeleteLocationRequest) (*emptypb.Empty, error) {
	if err := s.store.DeleteLocation(ctx, req.GetLocationName()); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to add location (%+v): %v", req.GetLocationName(), err)
	}
	return &emptypb.Empty{}, nil
}

// ListLocations returns all Locations.
func (s *Service) ListLocations(ctx context.Context, req *svpb.ListLocationsRequest) (*svpb.ListLocationsResponse, error) {
	locs, err := s.store.ListLocations(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list locations: %v", err)
	}
	resp := &svpb.ListLocationsResponse{}
	for _, l := range locs {
		resp.Locations = append(resp.Locations, resources.LocationToProto(l))
	}
	return resp, nil
}

// Process receives Fleetspeak messages and stores the enclosed SensorResponse.
func (s *Service) Process(ctx context.Context, m *fspb.Message) (*fspb.EmptyMessage, error) {
	var msg spb.SensorMessage
	if err := ptypes.UnmarshalAny(m.Data, &msg); err != nil {
		log.Errorf("Failed to unmarshal sensor message (%s)", m.Data.String())
		return &fspb.EmptyMessage{}, nil
	}
	log.Infof("Received sensor message (%X) from Fleetspeak", m.GetSourceMessageId())
	switch t := msg.Type.(type) {
	case *spb.SensorMessage_Response:
		req := resources.ProtoToSensorRequest(&msg)
		if err := s.store.ModifySensorRequest(ctx, req); err != nil {
			log.Errorf("Failed to update sensor request (%+v)", req)
		}
	case *spb.SensorMessage_Alert:
		if err := s.store.AddSensorMessage(ctx, resources.ProtoToSensorMessage(&msg)); err != nil {
			log.Errorf("Failed to store sensor alert (%+v)", msg.GetAlert())
		}
	case *spb.SensorMessage_Heartbeat:
		if err := s.store.AddSensorMessage(ctx, resources.ProtoToSensorMessage(&msg)); err != nil {
			log.Errorf("Failed to store sensor heartbeat (%+v)", msg.GetHeartbeat())
		}
	default:
		log.Errorf("Unknown sensor message type (%T)", t)
	}
	return &fspb.EmptyMessage{}, nil
}
