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

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/youtube/v3"
)

func TestResolveChannelTarget(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		mine    bool
		wantID  string
		wantErr bool
	}{
		{
			name:   "explicit_channel_id",
			args:   []string{"UC_x5XG1OV2P6uZZ5FSM9Ttw"},
			mine:   false,
			wantID: "UC_x5XG1OV2P6uZZ5FSM9Ttw",
		},
		{
			name:   "mine_flag_only",
			args:   nil,
			mine:   true,
			wantID: "",
		},
		{
			name:    "missing_both_id_and_mine",
			args:    nil,
			mine:    false,
			wantErr: true,
		},
		{
			name:    "conflicting_id_and_mine",
			args:    []string{"UC_x5XG1OV2P6uZZ5FSM9Ttw"},
			mine:    true,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveChannelTarget(tc.args, tc.mine)
			if tc.wantErr {
				if err == nil {
					t.Errorf("resolveChannelTarget(%v, %v) = %q, want error", tc.args, tc.mine, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveChannelTarget(%v, %v) returned unexpected error: %v", tc.args, tc.mine, err)
			}
			if got != tc.wantID {
				t.Errorf("resolveChannelTarget(%v, %v) = %q, want %q", tc.args, tc.mine, got, tc.wantID)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxRunes int
		want     string
	}{
		{
			name:     "short_string",
			input:    "hello",
			maxRunes: 10,
			want:     "hello",
		},
		{
			name:     "exact_length",
			input:    "hello",
			maxRunes: 5,
			want:     "hello",
		},
		{
			name:     "truncated_ascii",
			input:    "hello world",
			maxRunes: 5,
			want:     "hello...",
		},
		{
			name:     "truncated_multibyte_utf8",
			input:    "こんにちは世界",
			maxRunes: 5,
			want:     "こんにちは...",
		},
		{
			name:     "zero_limit",
			input:    "hello",
			maxRunes: 0,
			want:     "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := truncate(tc.input, tc.maxRunes); got != tc.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tc.input, tc.maxRunes, got, tc.want)
			}
		})
	}
}

func TestFilterPlaylistItems(t *testing.T) {
	pubItem := &youtube.PlaylistItem{
		Id:     "pub-1",
		Status: &youtube.PlaylistItemStatus{PrivacyStatus: "public"},
	}
	unlistedItem := &youtube.PlaylistItem{
		Id:     "unl-1",
		Status: &youtube.PlaylistItemStatus{PrivacyStatus: "unlisted"},
	}
	privItem := &youtube.PlaylistItem{
		Id:     "prv-1",
		Status: &youtube.PlaylistItemStatus{PrivacyStatus: "private"},
	}
	nilStatusItem := &youtube.PlaylistItem{
		Id: "nil-status",
	}
	allItems := []*youtube.PlaylistItem{pubItem, unlistedItem, privItem, nilStatusItem}

	tests := []struct {
		name      string
		filterCSV string
		wantIDs   []string
	}{
		{
			name:      "empty_filter_returns_all",
			filterCSV: "",
			wantIDs:   []string{"pub-1", "unl-1", "prv-1", "nil-status"},
		},
		{
			name:      "single_privacy_filter",
			filterCSV: "public",
			wantIDs:   []string{"pub-1"},
		},
		{
			name:      "multiple_privacy_statuses_with_spaces",
			filterCSV: " public , unlisted ",
			wantIDs:   []string{"pub-1", "unl-1"},
		},
		{
			name:      "no_matching_status",
			filterCSV: "unknown",
			wantIDs:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotItems := filterPlaylistItems(allItems, tc.filterCSV)
			var gotIDs []string
			for _, item := range gotItems {
				gotIDs = append(gotIDs, item.Id)
			}
			if diff := cmp.Diff(tc.wantIDs, gotIDs); diff != "" {
				t.Errorf("filterPlaylistItems(allItems, %q) IDs diff (-want +got):\n%s", tc.filterCSV, diff)
			}
		})
	}
}

func TestFormatDryRunUpload(t *testing.T) {
	t.Run("json_mode", func(t *testing.T) {
		var buf bytes.Buffer
		err := formatDryRunUpload(&buf, "./demo.mp4", "Demo Title", "Demo Desc", true)
		if err != nil {
			t.Fatalf("formatDryRunUpload(&buf, %q, %q, %q, true) returned unexpected error: %v", "./demo.mp4", "Demo Title", "Demo Desc", err)
		}

		var got map[string]any
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("json.Unmarshal(%q) returned unexpected error: %v", buf.String(), err)
		}

		want := map[string]any{
			"status": "dry-run",
			"file":   "./demo.mp4",
			"snippet": map[string]any{
				"title":       "Demo Title",
				"description": "Demo Desc",
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("formatDryRunUpload(&buf, %q, %q, %q, true) JSON diff (-want +got):\n%s", "./demo.mp4", "Demo Title", "Demo Desc", diff)
		}
	})

	t.Run("human_mode", func(t *testing.T) {
		var buf bytes.Buffer
		err := formatDryRunUpload(&buf, "./demo.mp4", "Demo Title", "Demo Desc", false)
		if err != nil {
			t.Fatalf("formatDryRunUpload(&buf, %q, %q, %q, false) returned unexpected error: %v", "./demo.mp4", "Demo Title", "Demo Desc", err)
		}

		got := buf.String()
		for _, wantSub := range []string{"DRY RUN MODE ENABLED", "./demo.mp4", "Demo Title", "Demo Desc"} {
			if !strings.Contains(got, wantSub) {
				t.Errorf("formatDryRunUpload(&buf, %q, %q, %q, false) = %q, want substring %q", "./demo.mp4", "Demo Title", "Demo Desc", got, wantSub)
			}
		}
	})
}

func TestCommandRegistrationAndRequiredFlags(t *testing.T) {
	subcommands := [][]string{
		{"channel"},
		{"channel", "info"},
		{"channel", "videos"},
		{"video"},
		{"video", "channel-info"},
		{"video", "info"},
		{"video", "upload"},
	}

	for _, path := range subcommands {
		name := strings.Join(path, "_")
		t.Run(name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find(path)
			if err != nil || cmd == nil {
				t.Fatalf("rootCmd.Find(%v) = (%v, %v), want registered command", path, cmd, err)
			}
			if cmd.Short == "" || cmd.Long == "" || cmd.Example == "" {
				t.Errorf("rootCmd.Find(%v) missing documentation pillars: Short=%q Long=%q Example=%q", path, cmd.Short, cmd.Long, cmd.Example)
			}
		})
	}

	t.Run("video_upload_requires_title", func(t *testing.T) {
		cmd, _, err := rootCmd.Find([]string{"video", "upload"})
		if err != nil {
			t.Fatalf("rootCmd.Find([\"video\", \"upload\"]) returned unexpected error: %v", err)
		}
		if err := cmd.Flags().Set("title", ""); err != nil {
			t.Fatalf("cmd.Flags().Set(\"title\", \"\") returned unexpected error: %v", err)
		}
		cmd.Flags().Lookup("title").Changed = false
		if err := cmd.ValidateRequiredFlags(); err == nil {
			t.Errorf("cmd.ValidateRequiredFlags() on %q = nil, want error for missing --title", cmd.CommandPath())
		}
	})

	t.Run("version_flag_outputs_version", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		rootCmd.SetArgs([]string{"--version"})
		t.Cleanup(func() {
			rootCmd.SetArgs(nil)
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("rootCmd.Execute() with --version returned unexpected error: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, Version) {
			t.Errorf("rootCmd --version output = %q, want substring %q", got, Version)
		}
	})
}
