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

// Package cmd defines the Cobra command hierarchy and CLI execution logic for yt-cli.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghchinoy/yt-cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const authSetupHint = "Place your downloaded OAuth client_secret.json in ~/.config/yt-cli/client_secret.json or set YT_CLIENT_ID and YT_CLIENT_SECRET."

var (
	cfgFile string
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:   "yt-cli",
	Short: "A CLI for interacting with YouTube API",
	Long:  `A Dual-Mode CLI for humans and AI agents to interact with YouTube API. Supports headless authentication and deterministic output.`,
	Example: `  yt-cli channel info --mine
  yt-cli channel videos --mine --limit 5 --json
  yt-cli video channel-info dQw4w9WgXcQ --json
  yt-cli video upload ./my-video.mp4 --title "My Title" --dry-run`,
}

// Execute runs the root yt-cli command and exits with a non-zero code on failure.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/yt-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output pure JSON for agents")
	cobra.CheckErr(viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json")))

	rootCmd.AddGroup(&cobra.Group{ID: "Info", Title: "Information Commands"})
	rootCmd.AddGroup(&cobra.Group{ID: "Upload", Title: "Upload Commands"})

	// Apply semantic coloring to help if interactive.
	if os.Getenv("NO_COLOR") == "" && os.Getenv("YT_CLI_NO_TUI") == "" {
		usageTemplate := rootCmd.UsageTemplate()

		coloredTemplate := strings.Replace(usageTemplate, "Usage:", ui.Accent.Render("Usage:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Information Commands", ui.Accent.Render("Information Commands:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Upload Commands", ui.Accent.Render("Upload Commands:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Additional Commands:", ui.Accent.Render("Additional Commands:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Flags:", ui.Accent.Render("Flags:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Available Commands:", ui.Accent.Render("Available Commands:"), 1)
		coloredTemplate = strings.Replace(coloredTemplate, "Global Flags:", ui.Accent.Render("Global Flags:"), 1)

		rootCmd.SetUsageTemplate(coloredTemplate)
	}
}

func initConfig() {
	configureViperPaths(cfgFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("YT")

	if err := viper.ReadInConfig(); err == nil && !viper.GetBool("json") {
		// Outputting to stderr prevents polluting stdout JSON streams.
		fmt.Fprintln(os.Stderr, ui.Muted.Render("Using config file:"), viper.ConfigFileUsed())
	}
}

func configureViperPaths(configFile string) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		return
	}

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configDir := filepath.Join(home, ".config", "yt-cli")
	cobra.CheckErr(os.MkdirAll(configDir, 0o700))

	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
}

func printJSON(w io.Writer, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON output: %w", err)
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}
