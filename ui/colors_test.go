// Copyright 2026 Google LLC
//
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

package ui

import (
	"errors"
	"strings"
	"testing"
)

func TestFormatError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		hint       string
		wantSubstr []string
		wantEmpty  bool
	}{
		{
			name:      "nil_error",
			err:       nil,
			hint:      "ignored hint",
			wantEmpty: true,
		},
		{
			name:       "error_without_hint",
			err:        errors.New("channel not found"),
			hint:       "",
			wantSubstr: []string{"Error: ", "channel not found"},
		},
		{
			name:       "error_with_hint",
			err:        errors.New("missing credentials"),
			hint:       "Place client_secret.json in ~/.config/yt-cli",
			wantSubstr: []string{"Error: ", "missing credentials", "Hint: ", "Place client_secret.json in ~/.config/yt-cli"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatError(tc.err, tc.hint)
			if tc.wantEmpty {
				if got != "" {
					t.Errorf("FormatError(%v, %q) = %q, want %q", tc.err, tc.hint, got, "")
				}
				return
			}
			for _, sub := range tc.wantSubstr {
				if !strings.Contains(got, sub) {
					t.Errorf("FormatError(%v, %q) = %q, want substring %q", tc.err, tc.hint, got, sub)
				}
			}
			if tc.hint == "" && strings.Contains(got, "Hint:") {
				t.Errorf("FormatError(%v, %q) = %q, want no Hint line when hint is empty", tc.err, tc.hint, got)
			}
		})
	}
}
