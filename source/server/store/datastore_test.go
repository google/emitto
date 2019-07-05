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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// This test relies on the Google Cloud Datastore emulator with the following configuration
// constants:
const (
	testHost    = "127.0.0.1:9999"
	testProject = "test-project-name"
)

func TestDataStore(t *testing.T) {
	RunTestSuite(t, func() (Store, error) {
		// Reset emulator before each test.
		if err := resetEmulator(); err != nil {
			return nil, err
		}
		c, err := datastore.NewClient(context.Background(), testProject, option.WithEndpoint(testHost), option.WithoutAuthentication(), option.WithGRPCDialOption(grpc.WithInsecure()))
		if err != nil {
			return nil, fmt.Errorf("GCD client creation failed: %v", err)
		}
		return &DataStore{c}, nil
	})
}

func resetEmulator() error {
	resp, err := http.Post(fmt.Sprintf("http://%s/reset", testHost), "", bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) != "Resetting...\n" {
		return errors.New("emulator failed to reset")
	}
	return nil
}
