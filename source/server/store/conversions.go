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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	"github.com/google/emitto/source/resources"
)

// TimeNow is stubbed for testing.
var TimeNow = time.Now

// MutateRule applies mutable, non-empty field mutations from the src to dst Rule.
func MutateRule(src, dst *resources.Rule) error {
	m, err := resources.MutationsMapping(resources.Rule{})
	if err != nil {
		return err
	}
	return mutateFields(reflect.ValueOf(*src), dst, m)
}

// MutateLocation applies mutable, non-empty field mutations from the src to dst Location.
func MutateLocation(src, dst *resources.Location) error {
	m, err := resources.MutationsMapping(resources.Location{})
	if err != nil {
		return err
	}
	return mutateFields(reflect.ValueOf(*src), dst, m)
}

// MutateSensorRequest applies mutable, non-empty field mutations from the src to dst SensorRequest.
func MutateSensorRequest(src, dst *resources.SensorRequest) error {
	m, err := resources.MutationsMapping(resources.SensorRequest{})
	if err != nil {
		return err
	}
	return mutateFields(reflect.ValueOf(*src), dst, m)
}

func mutateFields(src reflect.Value, dst interface{}, fields map[string]bool) error {
	for i := 0; i < src.NumField(); i++ {
		n := src.Type().Field(i).Name
		if fields[strings.ToLower(strings.Join(camelcase.Split(n), "_"))] {
			f := src.Field(i)
			if f.IsValid() && !isZero(f) {
				var d reflect.Value
				switch t := dst.(type) {
				case *resources.Rule:
					d = reflect.ValueOf(t).Elem()
				case *resources.Location:
					d = reflect.ValueOf(t).Elem()
				case *resources.SensorRequest:
					d = reflect.ValueOf(t).Elem()
				default:
					return fmt.Errorf("invalid mutable type: %T", t)
				}
				d.FieldByName(n).Set(f)
			}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice: // Currently supported types with default values of nil.
		return v.IsNil()
	}
	// All other types.
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
