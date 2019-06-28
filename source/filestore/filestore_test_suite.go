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
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"context"

	"github.com/google/go-cmp/cmp"

	log "github.com/golang/glog"
)

// suite contains store tests and the underlying FileStore implementation.
type suite struct {
	builder func() FileStore
}

// RunTestSuite runs all generic store tests.
//
// The tests use the provided builder to instantiate a FileStore.
// The builder is expected to always return a valid Store.
func RunTestSuite(t *testing.T, builder func() FileStore) {
	s := &suite{builder}
	s.Run(t)
}

// Run runs Test* methods of the suite as subtests.
func (s *suite) Run(t *testing.T) {
	sv := reflect.ValueOf(s)
	st := reflect.TypeOf(s)
	for i := 0; i < sv.NumMethod(); i++ {
		n := st.Method(i).Name
		if strings.HasPrefix(n, "Test") {
			mv := sv.MethodByName(n)
			mt := mv.Type()
			if mt.NumIn() != 1 || !reflect.TypeOf(t).AssignableTo(mt.In(0)) {
				log.Fatalf("Method %q of the test suite must have 1 argument of type *testing.T", n)
			}
			if mt.NumOut() != 0 {
				log.Fatalf("Method %q of the test suite must have no return value", n)
			}
			m := mv.Interface().(func(t *testing.T))
			t.Run(n, m)
		}
	}
}

var (
	dir1   = "/location1/zone1/"
	path1  = filepath.Join(dir1, "rulefile1")
	rules1 = []byte("rule 1\nrule 2\nrule3")

	path2  = filepath.Join(dir1, "rulefile2")
	rules2 = []byte("rule 4\nrule 5\nrule6")

	dir2   = "/location2/zone2/"
	path3  = filepath.Join(dir2, "rulefile1")
	rules3 = []byte("rule 7\nrule 8\nrule9")

	dir3   = "/location3/zone3/"
	path4  = filepath.Join(dir3, "rulefile1")
	rules4 = []byte("rule 1\nrule 2\nrule3")
)

func (s *suite) TestAddRuleFile(t *testing.T) {
	st := s.builder()
	ctx := context.Background()

	if err := st.AddRuleFile(ctx, path1, rules1); err != nil {
		t.Error(err)
	}
}

func (s *suite) TestGetRuleFile(t *testing.T) {
	st := s.builder()
	ctx := context.Background()

	if err := st.AddRuleFile(ctx, path1, rules1); err != nil {
		t.Error(err)
	}
	got, err := st.GetRuleFile(ctx, path1)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(rules1, got); diff != "" {
		t.Errorf("expectation mismatch:\n%s", diff)
	}
}

func (s *suite) TestDeleteRuleFile(t *testing.T) {
	st := s.builder()
	ctx := context.Background()

	if err := st.AddRuleFile(ctx, path1, rules1); err != nil {
		t.Error(err)
	}
	if err := st.DeleteRuleFile(ctx, path1); err != nil {
		t.Error(err)
	}
	if _, err := st.GetRuleFile(ctx, path1); err == nil {
		t.Error(errors.New("returning a non-existent rule file should have raised an error"))
	}
}
