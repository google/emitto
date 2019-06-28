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

package client

import (
	"fmt"
	"io"
	"testing"
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"

	pb "github.com/google/emitto/source/server/proto"
)

func TestGetSid(t *testing.T) {
	rule := `alert http any any -> any any (msg:"Test"; content:"Test"; nocase; classtype:policy-violation; sid:1234567890; rev:1;)`
	want := int64(1234567890)
	got, err := getSID(rule)
	if err != nil {
		t.Errorf("getSID() unexpected failure: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("getSID() expectation mismatch (-want +got):\n%s", diff)
	}
}

func TestDeployRules(t *testing.T) {
	responses := []*pb.DeployRulesResponse{
		{ClientId: "id1"},
		{ClientId: "id2"},
	}
	c := Client{emitto: &fakeEmittoClient{stream: &fakeEmittoDeployRulesClient{responses: responses}}}
	got, err := c.DeployRules(context.Background(), &pb.Location{})
	if err != nil {
		t.Errorf("DeployRules() retuned unexpected error: %v", err)
	}
	if len(responses) != len(got) {
		t.Errorf("DeployRules() expected: %v responses, got: %v", len(responses), len(got))
	}
	for i, r := range responses {
		if diff := cmp.Diff(r, got[i], cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("DeployRules() expectation mismatch (-want +got):\n%s", diff)
		}
	}

	wantErr := fmt.Errorf("error")
	c = Client{emitto: &fakeEmittoClient{stream: &fakeEmittoDeployRulesClient{err: wantErr}}}
	_, err = c.DeployRules(context.Background(), &pb.Location{})
	if diff := cmp.Diff(wantErr.Error(), err.Error()); diff != "" {
		t.Errorf("DeployRules() expectation mismatch (-want +got):\n%s", diff)
	}
}

// fakeEmittoClient fakes emitto client, field of Client struct object.
// This does not fake the implementation of DeployRules method of an actual
// Client struct, rather it only fakes the stream to the service.
type fakeEmittoClient struct {
	pb.EmittoClient
	stream *fakeEmittoDeployRulesClient
}

func (c *fakeEmittoClient) DeployRules(ctx context.Context, in *pb.DeployRulesRequest, opts ...grpc.CallOption) (pb.Emitto_DeployRulesClient, error) {
	return c.stream, nil
}

type fakeEmittoDeployRulesClient struct {
	grpc.ClientStream
	responses []*pb.DeployRulesResponse
	err       error
}

func (s *fakeEmittoDeployRulesClient) Recv() (*pb.DeployRulesResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	if len(s.responses) == 0 {
		return nil, io.EOF
	}
	r := s.responses[0]
	s.responses = s.responses[1:]
	return r, nil
}
