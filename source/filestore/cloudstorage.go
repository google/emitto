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

package filestore

import (
	"fmt"
	"io/ioutil"
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSFileStore is a Google Cloud Storage implementation of a FileStore.
type GCSFileStore struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

// NewGCSFileStore returns a new GCSFileStore.
func NewGCSFileStore(bucket string, client *storage.Client) *GCSFileStore {
	return &GCSFileStore{
		client: client,
		bucket: client.Bucket(bucket),
	}
}

// NewGCSClient initializes a new Google Cloud Storage Client.
// Follow these instructions to set up application credentials:
// https://cloud.google.com/docs/authentication/production#obtaining_and_providing_service_account_credentials_manually.
func NewGCSClient(ctx context.Context, credFile string, scopes []string) (*storage.Client, error) {
	c, err := storage.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, fmt.Errorf("GCS client creation failed: %v", err)
	}
	return c, nil
}

// AddRuleFile uploads a rule file to GCS.
func (s *GCSFileStore) AddRuleFile(ctx context.Context, path string, rules []byte) error {
	writer := s.bucket.Object(path).NewWriter(ctx)
	defer writer.Close()
	if _, err := writer.Write(rules); err != nil {
		return fmt.Errorf("writing to rule file %q failed: %v", path, err)
	}
	return nil
}

// GetRuleFile returns a rule file from GCS.
func (s *GCSFileStore) GetRuleFile(ctx context.Context, path string) ([]byte, error) {
	reader, err := s.bucket.Object(path).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("reader creation failed: %v", err)
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("rule file %q download failed: %v", path, err)
	}
	return data, nil
}

// DeleteRuleFile removes a rule file from GCS.
func (s *GCSFileStore) DeleteRuleFile(ctx context.Context, ruleFile string) error {
	obj := s.bucket.Object(ruleFile)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("object deletion failed: %v", err)
	}
	return nil
}

// Close GCS client connection.
func (s *GCSFileStore) Close() error {
	return s.client.Close()
}
