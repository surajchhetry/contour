// Copyright © 2018 Heptio
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contour

import (
	"math"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
)

func TestParseAnnotationTimeout(t *testing.T) {
	tests := map[string]struct {
		a    map[string]string
		want time.Duration
		ok   bool
	}{
		"nada": {
			a:    nil,
			want: 0,
			ok:   false,
		},
		"empty": {
			a:    map[string]string{annotationRequestTimeout: ""}, // not even sure this is possible via the API
			want: 0,
			ok:   false,
		},
		"infinity": {
			a:    map[string]string{annotationRequestTimeout: "infinity"},
			want: 0,
			ok:   true,
		},
		"10 seconds": {
			a:    map[string]string{annotationRequestTimeout: "10s"},
			want: 10 * time.Second,
			ok:   true,
		},
		"invalid": {
			a:    map[string]string{annotationRequestTimeout: "10"}, // 10 what?
			want: 0,
			ok:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := parseAnnotationTimeout(tc.a, annotationRequestTimeout)
			if got != tc.want || ok != tc.ok {
				t.Fatalf("parseAnnotationTimeout(%q): want: %v, %v, got: %v, %v", tc.a, tc.want, tc.ok, got, ok)
			}
		})
	}
}

func TestParseAnnotationUInt32(t *testing.T) {
	tests := map[string]struct {
		a     map[string]string
		want  uint32
		isNil bool
	}{
		"nada": {
			a:     nil,
			isNil: true,
		},
		"empty": {
			a:     map[string]string{annotationRequestTimeout: ""}, // not even sure this is possible via the API
			isNil: true,
		},
		"smallest": {
			a:    map[string]string{annotationRequestTimeout: "0"},
			want: 0,
		},
		"middle value": {
			a:    map[string]string{annotationRequestTimeout: "20"},
			want: 20,
		},
		"biggest": {
			a:    map[string]string{annotationRequestTimeout: "4294967295"},
			want: math.MaxUint32,
		},
		"invalid": {
			a:     map[string]string{annotationRequestTimeout: "10seconds"}, // not a duration
			isNil: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseAnnotationUInt32(tc.a, annotationRequestTimeout)
			full := types.UInt32Value{Value: tc.want}

			if ((got == nil) != tc.isNil) || (got != nil && *got != full) {
				t.Fatalf("parseAnnotationUInt32(%q): want: %v, isNil: %v, got: %v", tc.a, tc.want, tc.isNil, got)
			}
		})
	}
}
