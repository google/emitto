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

package fleetspeak

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/fleetspeak/fleetspeak/src/client/channel"
	"github.com/google/fleetspeak/fleetspeak/src/client/service"
	"github.com/google/go-cmp/cmp"

	log "github.com/golang/glog"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
)

func TestSendAndWait(t *testing.T) {
	o := make(chan service.AckMessage)
	callbackChan := make(chan string)
	defer close(o)

	// Client with fake relentless channel to Fleetspeak client.
	c := &Client{
		fsChan: &channel.RelentlessChannel{
			Out: o,
		},
		callbackChan: callbackChan,
	}

	// Fleetspeak client acknowledgement message.
	ackMsg := &service.AckMessage{
		Ack: func() {
			c.callbackChan <- "TEST_OP_ID"
		},
	}
	go c.sendAndWait(ackMsg)

	// Fake Fleetspeak client message handling.
	go func() {
		for {
			select {
			case <-o:
				log.Info("TestSendAndWait() Fleetspeak client received message")
				time.Sleep(5 * time.Second) // Simulate work.
				callbackChan <- "TEST_OP_ID"
				return
			default:
				log.Info("No new messages")
			}
		}
	}()

	want := "TEST_OP_ID"
	got := <-c.callbackChan
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestSendAndWait() expectation mismatch:\n%s", diff)
	}
}

func TestReceive(t *testing.T) {
	in := make(chan *fspb.Message)

	// Client with fake relentless channel to Fleetspeak client.
	c := &Client{
		fsChan: &channel.RelentlessChannel{
			In: in,
		},
		messages: make(chan *fspb.Message, maxMessages),
	}

	// Fake Fleetspeak client message sending.
	// Send 3 messages.
	m := &fspb.Message{
		MessageId: []byte{1, 2, 3},
		Data: &any.Any{
			Value: []byte{1, 2, 3},
		},
	}

	go func() {
		for i := 0; i < 3; i++ {
			in <- m
		}
	}()

	// Receive messages.
	done := make(chan struct{})
	go func() {
		c.Receive(done)
	}()

	// Simulate work, then close receiving loop.
	time.Sleep(5 * time.Second)
	close(done)

	// Verify received messages.
	want := &fspb.Message{
		MessageId: []byte{1, 2, 3},
		Data: &any.Any{
			Value: []byte{1, 2, 3},
		},
	}

	for got := range c.Messages() {
		if diff := cmp.Diff(want, got, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("TestReceive() expectation mismatch (-want +got):\n%s", diff)
		}
	}
}
