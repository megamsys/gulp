/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package toml_test

import (
	"reflect"
	"testing"

	"github.com/megamsys/gulp/toml"
)

// Ensure that megabyte sizes can be parsed.
func TestSize_UnmarshalText_MB(t *testing.T) {
	var s toml.Size
	if err := s.UnmarshalText([]byte("200m")); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if s != 200*(1<<20) {
		t.Fatalf("unexpected size: %d", s)
	}
}

// Ensure that gigabyte sizes can be parsed.
func TestSize_UnmarshalText_GB(t *testing.T) {
	if typ := reflect.TypeOf(0); typ.Size() != 8 {
		t.Skip("large gigabyte parsing on 64-bit arch only")
	}

	var s toml.Size
	if err := s.UnmarshalText([]byte("10g")); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else if s != 10*(1<<30) {
		t.Fatalf("unexpected size: %d", s)
	}
}
