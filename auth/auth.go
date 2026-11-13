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

// Package auth manages OAuth 2.0 credentials and initializes the YouTube Data API v3 client.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// NewYouTubeService creates an authenticated YouTube Data API v3 service client.
func NewYouTubeService(ctx context.Context) (*youtube.Service, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "yt-cli")
	secretFile := filepath.Join(configDir, "client_secret.json")
	tokenFile := filepath.Join(configDir, "token.json")

	config, err := loadOAuthConfig(secretFile)
	if err != nil {
		return nil, err
	}

	client, err := newHTTPClient(ctx, config, tokenFile)
	if err != nil {
		return nil, err
	}

	srv, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("creating YouTube client: %w", err)
	}
	return srv, nil
}

func loadOAuthConfig(secretFile string) (*oauth2.Config, error) {
	b, err := os.ReadFile(secretFile)
	if err != nil {
		clientID := os.Getenv("YT_CLIENT_ID")
		clientSecret := os.Getenv("YT_CLIENT_SECRET")
		if clientID != "" && clientSecret != "" {
			return &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     google.Endpoint,
				RedirectURL:  "http://localhost",
				Scopes:       []string{youtube.YoutubeScope},
			}, nil
		}
		return nil, fmt.Errorf("reading client secret file %q: %w", secretFile, err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		return nil, fmt.Errorf("parsing client secret file %q: %w", secretFile, err)
	}
	return config, nil
}

func newHTTPClient(ctx context.Context, config *oauth2.Config, tokenFile string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = fetchTokenFromWeb(ctx, config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokenFile, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

func fetchTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Fprintf(os.Stderr, "Go to the following link in your browser then type the authorization code:\n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("reading authorization code: %w", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("exchanging authorization code for token: %w", err)
	}
	return tok, nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening token file %q: %w", path, err)
	}
	defer f.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, fmt.Errorf("decoding token file %q: %w", path, err)
	}
	return tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating token directory for %q: %w", path, err)
	}
	fmt.Fprintf(os.Stderr, "Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("opening token file %q for writing: %w", path, err)
	}
	encodeErr := json.NewEncoder(f).Encode(token)
	closeErr := f.Close()
	if encodeErr != nil {
		return fmt.Errorf("encoding oauth token to %q: %w", path, encodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("closing token file %q: %w", path, closeErr)
	}
	return nil
}
