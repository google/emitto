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

// Package fleetspeak provides functionality for network sensors to communicate with the Emitto
// service via Fleetspeak.
package fleetspeak

import (
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/fleetspeak/fleetspeak/src/client/channel"
	"github.com/google/fleetspeak/fleetspeak/src/client/service"
	"github.com/google/fleetspeak/fleetspeak/src/client/socketservice/client"

	log "github.com/golang/glog"
	pb "github.com/google/emitto/source/sensor/proto"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
)

const (
	// Service name for Fleetspeak messages.
	serviceName = "Emitto"
	// Maximum size of Messages channel used to receive Fleetspeak client messages.
	maxMessages = 1
)

// Client contains functionality to send and receive messages to/from a Fleetspeak Client.
type Client struct {
	// Channel used for sending messages to the Fleetspeak client.
	fsChan *channel.RelentlessChannel
	// Callback for Fleetspeak client send acknowledgements.
	callbackChan chan string
	// Messages is used to queue received Fleetspeak client messages for sensor client consumption.
	messages chan *fspb.Message
}

// New initializes a Client.
func New(socket string) *Client {
	rc := client.OpenChannel(socket, time.Now().Format(time.RFC1123Z))

	return &Client{
		fsChan:       rc,
		callbackChan: make(chan string, 5), // To prevent potential locking.
		messages:     make(chan *fspb.Message, maxMessages),
	}
}

// SendMessage a message to the Fleetspeak client. This call blocks until Fleetspeak has
// acknowledged the message.
func (c *Client) SendMessage(m *pb.SensorMessage) (string, error) {
	req, err := c.createRequest(m)
	if err != nil {
		return "", err
	}
	return c.sendAndWait(req), nil
}

// sendAndWait sends a message to the Fleetspeak client and waits indefinitely. Only one sendAndWait
// should be called by the Client at a time to avoid non-chronological request ID logging from the
// Fleetspeak client callback channel.
func (c *Client) sendAndWait(msg *service.AckMessage) string {
	c.fsChan.Out <- *msg
	log.Infof("Sent message (%X) to Fleetspeak; awaiting acknowledgement...", msg.M.GetSourceMessageId())
	ack := <-c.callbackChan
	log.Infof("Received ack %q from Fleetspeak", ack)
	return ack
}

// createRequest composes a Fleetspeak AckMessage.
// Fleetspeak is optimized to handle messages sizes < 2MB.
func (c *Client) createRequest(m *pb.SensorMessage) (*service.AckMessage, error) {
	data, err := ptypes.MarshalAny(m)
	if err != nil {
		return nil, err
	}
	id := make([]byte, 16)
	rand.Read(id)
	return &service.AckMessage{
		M: &fspb.Message{
			SourceMessageId: id,
			Destination: &fspb.Address{
				ServiceName: serviceName,
			},
			Data:       data,
			Background: true,
		},
		Ack: func() {
			c.callbackChan <- m.Id
		},
	}, nil
}

// Receive continuously receives new messages from the Fleetspeak client's In channel. Once it
// receives a message, it will send it to the Messages channel for the sensor client to process.
func (c *Client) Receive(done <-chan struct{}) {
	for {
		select {
		case m := <-c.fsChan.In:
			log.Infof("Received message (%X) from Fleetspeak", m.GetSourceMessageId())
			c.messages <- m
		case <-done:
			log.Warning("Stopped receiving messages from Fleetspeak")
			close(c.messages)
			return
		}
	}
}

// Messages returns the channel containing incoming Fleetspeak messages.
func (c *Client) Messages() chan *fspb.Message {
	return c.messages
}
