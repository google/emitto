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
	"errors"
	"io"
	"net"
	"testing"
	"time"
	"context"

	"github.com/google/emitto/source/filestore"
	"github.com/google/emitto/source/resources"
	"github.com/google/emitto/source/server/fleetspeak"
	"github.com/google/emitto/source/server/store"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/golang/glog"
	spb "github.com/google/emitto/source/server/proto"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
	fsspb "github.com/google/fleetspeak/fleetspeak/src/server/proto/fleetspeak_server"
	mpb "google.golang.org/genproto/protobuf/field_mask"
)

func TestDeployRules(t *testing.T) {
	ctx := context.Background()
	timeNow = func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}

	tests := []struct {
		desc     string
		req      *spb.DeployRulesRequest
		fsServer *fakeFSAdminServer
		want     []*spb.DeployRulesResponse
		wantErr  bool
	}{
		{
			desc: "successful deployment",
			req:  &spb.DeployRulesRequest{Location: &spb.Location{Name: "a", Zones: []string{"dmz"}}},
			fsServer: &fakeFSAdminServer{
				listClients: func(*fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error) {
					return &fsspb.ListClientsResponse{Clients: testClients}, nil
				},
				insertMessage: func(*fspb.Message) (*fspb.EmptyMessage, error) {
					return &fspb.EmptyMessage{}, nil
				},
			},
			want: []*spb.DeployRulesResponse{
				{
					ClientId: "636C69656E745F61", // "client_a"
					Status:   status.New(codes.OK, "OK").Proto(),
				},
				{
					ClientId: "636C69656E745F62", // "client_b"
					Status:   status.New(codes.OK, "OK").Proto(),
				},
			},
		},
		{
			desc: "list client error",
			req:  &spb.DeployRulesRequest{Location: &spb.Location{Name: "a", Zones: []string{"dmz"}}},
			fsServer: &fakeFSAdminServer{
				listClients: func(*fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error) {
					return nil, errors.New("error")
				},
			},
			wantErr: true,
		},
		{
			desc: "insert message error",
			req:  &spb.DeployRulesRequest{Location: &spb.Location{Name: "a", Zones: []string{"dmz"}}},
			fsServer: &fakeFSAdminServer{
				listClients: func(*fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error) {
					return &fsspb.ListClientsResponse{Clients: testClients}, nil
				},
				insertMessage: func(*fspb.Message) (*fspb.EmptyMessage, error) {
					return nil, errors.New("error")
				},
			},
			want: []*spb.DeployRulesResponse{
				{
					ClientId: "636C69656E745F61", // "client_a"
					Status:   status.New(codes.Internal, `failed to insert message: failed to insert message for client (636C69656E745F61): rpc error: code = Unknown desc = error`).Proto(),
				},
				{
					ClientId: "636C69656E745F62", // "client_b"
					Status:   status.New(codes.Internal, `failed to insert message: failed to insert message for client (636C69656E745F62): rpc error: code = Unknown desc = error`).Proto(),
				},
			},
		},
		{
			desc: "no rules for location",
			req:  &spb.DeployRulesRequest{Location: &spb.Location{Name: "unknown", Zones: []string{"dmz"}}},
			fsServer: &fakeFSAdminServer{
				listClients: func(*fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error) {
					return &fsspb.ListClientsResponse{Clients: testClients}, nil
				},
				insertMessage: func(*fspb.Message) (*fspb.EmptyMessage, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// Set up storage.
		ds := store.NewMemoryStore()
		for _, r := range testRules {
			if err := ds.AddRule(ctx, r); err != nil {
				t.Fatal(err)
			}
		}
		fs := filestore.NewMemoryFileStore()
		// Set up servers and test clients.
		fc, stopFs := initFSAdminServerAndClient(t, tt.fsServer)
		defer fc.Close()
		defer stopFs()
		s := &Service{
			store:      ds,
			fileStore:  fs,
			fleetspeak: fc,
		}
		c, stopServer := initServerAndClient(t, s)
		defer stopServer()
		t.Run(tt.desc, func(t *testing.T) {
			stream, err := c.DeployRules(ctx, tt.req)
			if err != nil {
				t.Fatalf("failed to deploy rules: %v", err)
			}
			var got []*spb.DeployRulesResponse
			for {
				r, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("got err=%v, wantErr=%t", err, tt.wantErr)
				}
				if err != nil {
					return
				}
				got = append(got, r)
			}
			if diff := cmp.Diff(tt.want, got, cmp.Comparer(proto.Equal)); diff != "" {
				t.Errorf("expectation mismatch (want -> got):\n%s", diff)
			}
		})
	}
}

func TestModifyRule(t *testing.T) {
	ctx := context.Background()
	timeNow := func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}
	store.TimeNow = timeNow

	tests := []struct {
		desc    string
		rule    *spb.Rule
		mask    *mpb.FieldMask
		want    *resources.Rule
		wantErr bool
	}{
		{
			desc: "successfully modified",
			rule: &spb.Rule{
				Id:            1111,
				Body:          "updated",
				LocationZones: []string{"a:dmz", "b:corp"},
			},
			mask: &mpb.FieldMask{Paths: []string{"body"}},
			want: &resources.Rule{
				ID:           1111,
				Body:         "updated",
				LocZones:     []string{"a:dmz", "b:corp"},
				LastModified: timeNow().Format(time.RFC1123Z),
			},
		},
		{
			desc: "invalid mask path",
			rule: &spb.Rule{
				Id:            1111,
				Body:          "updated",
				LocationZones: []string{"a:dmz", "b:corp"},
			},
			mask:    &mpb.FieldMask{Paths: []string{"id"}},
			wantErr: true,
		},
	}
	// Set up storage.
	ds := store.NewMemoryStore()
	for _, r := range testRules {
		if err := ds.AddRule(ctx, r); err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range tests {
		// Set up servers and test clients.
		s := &Service{store: ds}
		c, stopServer := initServerAndClient(t, s)
		defer stopServer()
		t.Run(tt.desc, func(t *testing.T) {
			_, err := c.ModifyRule(ctx, &spb.ModifyRuleRequest{Rule: tt.rule, FieldMask: tt.mask})
			if (err != nil) != tt.wantErr {
				t.Errorf("got err=%v, wantErr=%t", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			got, err := s.store.ListRules(ctx, []int64{tt.want.ID})
			if err != nil {
				t.Error(err)
			}
			if l := len(got); l != 1 {
				t.Errorf("expected 1 rule, got %d", l)
			}
			if diff := cmp.Diff(tt.want, got[0]); diff != "" {
				t.Errorf("expectation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestModifyLocation(t *testing.T) {
	ctx := context.Background()
	tn := func() time.Time {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	}
	timeNow = tn
	store.TimeNow = tn

	tests := []struct {
		desc     string
		location *spb.Location
		mask     *mpb.FieldMask
		want     *resources.Location
		wantErr  bool
	}{
		{
			desc: "successfully modified",
			location: &spb.Location{
				Name:  "a",
				Zones: []string{"updated"},
			},
			mask: &mpb.FieldMask{Paths: []string{"zones"}},
			want: &resources.Location{
				Name:         "a",
				Zones:        []string{"updated"},
				LastModified: timeNow().Format(time.RFC1123Z),
			},
		},
		{
			desc: "invalid mask path",
			location: &spb.Location{
				Name:  "updated",
				Zones: []string{"updated"},
			},
			mask:    &mpb.FieldMask{Paths: []string{"name", "zones"}},
			wantErr: true,
		},
	}
	// Set up storage.
	ds := store.NewMemoryStore()
	for _, l := range testLocations {
		if err := ds.AddLocation(ctx, l); err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range tests {
		// Set up servers and testClients.
		s := &Service{store: ds}
		c, stopServer := initServerAndClient(t, s)
		defer stopServer()
		t.Run(tt.desc, func(t *testing.T) {
			_, err := c.ModifyLocation(ctx, &spb.ModifyLocationRequest{Location: tt.location, FieldMask: tt.mask})
			if (err != nil) != tt.wantErr {
				t.Errorf("got err=%v, wantErr=%t", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			got, err := s.store.GetLocation(ctx, tt.want.Name)
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("expectation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func initServerAndClient(t *testing.T, s *Service) (spb.EmittoClient, func()) {
	l, err := net.Listen("tcp", "localhost:")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	spb.RegisterEmittoServer(srv, s)
	go srv.Serve(l)

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	return spb.NewEmittoClient(conn), srv.Stop
}

func initFSAdminServerAndClient(t *testing.T, f *fakeFSAdminServer) (*fleetspeak.Client, func()) {
	l, err := net.Listen("tcp", "localhost:")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	fsspb.RegisterAdminServer(srv, f)
	go srv.Serve(l)

	fs, err := fleetspeak.New(l.Addr().String(), "")
	if err != nil {
		t.Fatal(err)
	}

	return fs, srv.Stop
}

type fakeFSAdminServer struct {
	fsspb.AdminServer

	insertMessage func(*fspb.Message) (*fspb.EmptyMessage, error)
	listClients   func(*fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error)
}

func (s *fakeFSAdminServer) InsertMessage(_ context.Context, m *fspb.Message) (*fspb.EmptyMessage, error) {
	return s.insertMessage(m)
}

func (s *fakeFSAdminServer) ListClients(_ context.Context, req *fsspb.ListClientsRequest) (*fsspb.ListClientsResponse, error) {
	return s.listClients(req)
}
