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

// Package filestore contains functionality to store Emitto rule files.
package filestore

import (
	"context"
)

// FileStore represents a Emitto file store.
type FileStore interface {
	// AddRuleFile creates a new rule file at the specified path.
	AddRuleFile(ctx context.Context, path string, rules []byte) error
	// GetRuleFile retrieves an existing rule file by path.
	GetRuleFile(ctx context.Context, path string) ([]byte, error)
	// DeleteRuleFile removes an existing rule file by path.
	DeleteRuleFile(ctx context.Context, path string) error
}
