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
	"os"
	"context"

	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MemoryFileStore a memory implementation of a FileStore.
type MemoryFileStore struct {
	store afero.Fs
}

// NewMemoryFileStore returns a MemoryFileStore.
func NewMemoryFileStore() *MemoryFileStore {
	return &MemoryFileStore{
		store: afero.NewMemMapFs(),
	}
}

// AddRuleFile stores a rule file in the filestore.
func (s *MemoryFileStore) AddRuleFile(ctx context.Context, path string, rules []byte) error {
	f, err := s.store.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to open rule file %q: %v", path, err))
	}
	defer f.Close()

	if _, err = f.Write(rules); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to write rules to %q: %v", path, err))
	}
	return nil
}

// GetRuleFile retrieves a rule file from the filestore.
func (s *MemoryFileStore) GetRuleFile(ctx context.Context, path string) ([]byte, error) {
	return afero.ReadFile(s.store, path)
}

// DeleteRuleFile removes a rule file from the filestore.
func (s *MemoryFileStore) DeleteRuleFile(ctx context.Context, path string) error {
	return s.store.Remove(path)
}
