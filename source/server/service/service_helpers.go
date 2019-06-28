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

package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/emitto/source/resources"

	spb "github.com/google/emitto/source/server/proto"
	fsspb "github.com/google/fleetspeak/fleetspeak/src/server/proto/fleetspeak_server"
	mpb "google.golang.org/genproto/protobuf/field_mask"
)

var timeNow = time.Now // Stubbed out for testing.

func ruleFilepath(location string) string {
	return filepath.Join(location, fmt.Sprintf("%s/%d", timeNow().Format("2006/01/02"), timeNow().Unix()))
}

func filterRulesByLocation(rules []*resources.Rule, loc *spb.Location) []*resources.Rule {
	matches := make(map[int64]*resources.Rule)
	for _, r := range rules {
		for _, lz := range r.LocZones {
			// fmt.Printf("LOCZONE: %+v\n", lz)
			// Filter by location name first.
			if strings.HasPrefix(lz, loc.GetName()+":") {
				for _, z := range loc.GetZones() {
					// Filter by zone name.
					if strings.HasSuffix(lz, ":"+z) {
						matches[r.ID] = r
					}
				}
			}
		}
	}
	res := make([]*resources.Rule, 0, len(matches))
	for _, m := range matches {
		res = append(res, m)
	}
	return res
}

func getClientIDsByLocation(clients []*fsspb.Client, loc *spb.Location) [][]byte {
	// Filter by location name.
	var filtered []*fsspb.Client
	for _, c := range clients {
		labels := c.GetLabels()
		for _, l := range labels {
			if l.GetLabel() == (resources.LocationNamePrefix + loc.GetName()) {
				filtered = append(filtered, c)
				break
			}
		}
	}
	// Filter by location zones.
	var ids [][]byte
	for _, c := range filtered {
		labels := c.GetLabels()
		var found bool
		for _, l := range labels {
			for _, z := range loc.GetZones() {
				if l.GetLabel() == (resources.LocationZonePrefix + z) {
					ids = append(ids, c.GetClientId())
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	return ids
}

// ValidateUpdateMask confirms whether the specified update field mask is valid for the
// respective object. If not, the invalid fields are returned in the error.
func ValidateUpdateMask(obj interface{}, mask *mpb.FieldMask) error {
	m, err := resources.MutationsMapping(obj)
	if err != nil {
		return err
	}
	var invalid []string
	for _, p := range mask.GetPaths() {
		if !m[p] {
			invalid = append(invalid, p)
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("the following fields are not mutable: %v", invalid)
	}
	return nil
}
