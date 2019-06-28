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

// Package fleetspeak provides administrative functionality for Fleetspeak.
package fleetspeak

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/fleetspeak/fleetspeak/src/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/golang/glog"
	spb "github.com/google/emitto/source/sensor/proto"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
	fsspb "github.com/google/fleetspeak/fleetspeak/src/server/proto/fleetspeak_server"
)

const (
	// Service name for Fleetspeak messages.
	serviceName = "Emitto"
	dateFmt     = "15:04:05.000 2006.01.02"
)

// Client represents an Fleetspeak admin client.
type Client struct {
	admin fsspb.AdminClient
	conn  *grpc.ClientConn
}

// New creates a new Client.
func New(addr, certFile string) (*Client, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	switch {
	case certFile != "":
		cred, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return nil, err
		}
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(cred))
	default:
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
	}
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Fleetspeak admin interface [%s]: %v", addr, err)
	}

	return &Client{
		admin: fsspb.NewAdminClient(conn),
		conn:  conn,
	}, nil
}

// Close terminates the Fleetspeak admin connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// InsertMessage inserts a message into the Fleetspeak system to be delivered to a sensor, where
// the sensor is identified by the id.
func (c *Client) InsertMessage(ctx context.Context, req *spb.SensorRequest, id []byte) error {
	data, err := ptypes.MarshalAny(req)
	if err != nil {
		return fmt.Errorf("failed to marshal operation message (%+v): %v", req, err)
	}
	mid := make([]byte, 16)
	rand.Read(mid)
	m := &fspb.Message{
		SourceMessageId: mid,
		Source: &fspb.Address{
			ServiceName: serviceName,
		},
		Destination: &fspb.Address{
			ServiceName: serviceName,
			ClientId:    id,
		},
		Data:       data,
		Background: true,
	}
	if _, err := c.admin.InsertMessage(ctx, m); err != nil {
		return fmt.Errorf("failed to insert message for client (%X): %v", id, err)
	}
	log.Infof("Sent message (%X) to Fleetspeak", m.GetSourceMessageId())
	return nil
}

// ListClients returns a list of clients from Fleetspeak.
func (c *Client) ListClients(ctx context.Context) ([]*fsspb.Client, error) {
	res, err := c.admin.ListClients(ctx, &fsspb.ListClientsRequest{}, grpc.MaxCallRecvMsgSize(1024*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %v", err)
	}
	if len(res.Clients) == 0 {
		return nil, errors.New("no clients found")
	}
	sort.Sort(byContactTime(res.Clients))
	return res.Clients, nil
}

// ParseClients returns human-readable client details.
func ParseClients(clients []*fsspb.Client) []string {
	var res []string
	for _, c := range clients {
		id, err := common.BytesToClientID(c.ClientId)
		if err != nil {
			log.Errorf("Ignoring invalid client (id=%v), %v", c.ClientId, err)
			continue
		}
		var ls []string
		for _, l := range c.Labels {
			ls = append(ls, l.ServiceName+":"+l.Label)
		}
		ts, err := ptypes.Timestamp(c.LastContactTime)
		if err != nil {
			log.Errorf("Unable to parse last contact time for client (id=%v): %v", id, err)
		}
		tag := ""
		if c.Blacklisted {
			tag = " *blacklisted*"
		}
		res = append(res, fmt.Sprintf("%v %v [%v]%s\n", id, ts.Format(dateFmt), strings.Join(ls, ","), tag))
	}
	return res
}

// byContactTime adapts []*fspb.Client for use by sort.Sort.
type byContactTime []*fsspb.Client

func (b byContactTime) Len() int           { return len(b) }
func (b byContactTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byContactTime) Less(i, j int) bool { return contactTime(b[i]).Before(contactTime(b[j])) }

func contactTime(c *fsspb.Client) time.Time {
	return time.Unix(c.LastContactTime.Seconds, int64(c.LastContactTime.Nanos))
}
