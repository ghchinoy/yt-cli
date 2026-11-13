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

package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func TestSaveAndLoadToken(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "nested", "token.json")
	want := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC),
	}

	if err := saveToken(tokenPath, want); err != nil {
		t.Fatalf("saveToken(%q, %+v) returned unexpected error: %v", tokenPath, want, err)
	}

	got, err := tokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("tokenFromFile(%q) returned unexpected error: %v", tokenPath, err)
	}

	if diff := cmp.Diff(want.AccessToken, got.AccessToken); diff != "" {
		t.Errorf("tokenFromFile(%q) AccessToken diff (-want +got):\n%s", tokenPath, diff)
	}
	if diff := cmp.Diff(want.RefreshToken, got.RefreshToken); diff != "" {
		t.Errorf("tokenFromFile(%q) RefreshToken diff (-want +got):\n%s", tokenPath, diff)
	}
}

func TestTokenFromFile_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	invalidJSONPath := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(invalidJSONPath, []byte("{not-json"), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) returned unexpected error: %v", invalidJSONPath, err)
	}

	tests := []struct {
		name string
		path string
	}{
		{
			name: "missing_file",
			path: filepath.Join(tmpDir, "nonexistent.json"),
		},
		{
			name: "malformed_json",
			path: invalidJSONPath,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tokenFromFile(tc.path)
			if err == nil {
				t.Errorf("tokenFromFile(%q) = %+v, want non-nil error", tc.path, got)
			}
		})
	}
}

func TestLoadOAuthConfig(t *testing.T) {
	tmpDir := t.TempDir()
	validSecretFile := filepath.Join(tmpDir, "client_secret.json")
	validSecretJSON := `{"installed":{"client_id":"file-client-id","client_secret":"file-client-secret","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","redirect_uris":["http://localhost"]}}`
	if err := os.WriteFile(validSecretFile, []byte(validSecretJSON), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) returned unexpected error: %v", validSecretFile, err)
	}

	invalidSecretFile := filepath.Join(tmpDir, "invalid_secret.json")
	if err := os.WriteFile(invalidSecretFile, []byte("not-valid-oauth-json"), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) returned unexpected error: %v", invalidSecretFile, err)
	}

	tests := []struct {
		name         string
		secretFile   string
		clientID     string
		clientSecret string
		wantClientID string
		wantErr      bool
	}{
		{
			name:         "from_secret_file",
			secretFile:   validSecretFile,
			wantClientID: "file-client-id",
		},
		{
			name:         "fallback_to_env_vars",
			secretFile:   filepath.Join(tmpDir, "missing.json"),
			clientID:     "env-client-id",
			clientSecret: "env-client-secret",
			wantClientID: "env-client-id",
		},
		{
			name:       "missing_file_and_env",
			secretFile: filepath.Join(tmpDir, "missing.json"),
			wantErr:    true,
		},
		{
			name:       "malformed_secret_file",
			secretFile: invalidSecretFile,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("YT_CLIENT_ID", tc.clientID)
			t.Setenv("YT_CLIENT_SECRET", tc.clientSecret)

			got, err := loadOAuthConfig(tc.secretFile)
			if tc.wantErr {
				if err == nil {
					t.Errorf("loadOAuthConfig(%q) = %+v, want error", tc.secretFile, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("loadOAuthConfig(%q) returned unexpected error: %v", tc.secretFile, err)
			}
			if got.ClientID != tc.wantClientID {
				t.Errorf("loadOAuthConfig(%q).ClientID = %q, want %q", tc.secretFile, got.ClientID, tc.wantClientID)
			}
		})
	}
}

func TestNewYouTubeService(t *testing.T) {
	ctx := t.Context()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configDir := filepath.Join(tmpHome, ".config", "yt-cli")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("os.MkdirAll(%q) returned unexpected error: %v", configDir, err)
	}

	secretFile := filepath.Join(configDir, "client_secret.json")
	secretJSON := `{"installed":{"client_id":"test-id","client_secret":"test-secret","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","redirect_uris":["http://localhost"]}}`
	if err := os.WriteFile(secretFile, []byte(secretJSON), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) returned unexpected error: %v", secretFile, err)
	}

	tokenFile := filepath.Join(configDir, "token.json")
	tok := &oauth2.Token{
		AccessToken: "cached-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}
	if err := saveToken(tokenFile, tok); err != nil {
		t.Fatalf("saveToken(%q, %+v) returned unexpected error: %v", tokenFile, tok, err)
	}

	srv, err := NewYouTubeService(ctx)
	if err != nil {
		t.Fatalf("NewYouTubeService(ctx) returned unexpected error: %v", err)
	}
	if srv == nil {
		t.Errorf("NewYouTubeService(ctx) = nil, want non-nil *youtube.Service")
	}
}
