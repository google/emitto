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

// Package client contains Emitto client functionality.
package client

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"context"

	"google.golang.org/grpc"

	log "github.com/golang/glog"
	pb "github.com/google/emitto/source/server/proto"
)

var suricataSIDRE = regexp.MustCompile(`sid:(\d+);`)

// Client represents a Emitto client.
type Client struct {
	conn   *grpc.ClientConn
	emitto pb.EmittoClient
}

// New returns new Client.
func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithDefaultCallOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Emitto (%s): %v", addr, err)
	}
	return &Client{
		conn:   conn,
		emitto: pb.NewEmittoClient(conn),
	}, nil
}

// Close Emitto client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// DeployRules deploys the rules to the provided location.
func (c *Client) DeployRules(ctx context.Context, loc *pb.Location) ([]*pb.DeployRulesResponse, error) {
	stream, err := c.emitto.DeployRules(ctx, &pb.DeployRulesRequest{Location: loc})
	if err != nil {
		return nil, err
	}

	var responses []*pb.DeployRulesResponse
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("deployment failure: %v", err)
			return responses, err
		}
		log.Infof("deployment response: %+v", resp)
		responses = append(responses, resp)
	}
	return responses, nil
}

// getSID extracts the SID from a Suricata rule and casts it to an int64.
func getSID(rule string) (int64, error) {
	matches := suricataSIDRE.FindStringSubmatch(rule)
	if len(matches) != 2 {
		return 0, errors.New("unable to properly locate rule SID")
	}
	sid, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to cast SID: %v", err)
	}
	return sid, nil
}
