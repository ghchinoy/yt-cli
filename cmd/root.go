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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ghchinoy/yt-cli/ui"
)

var (
	cfgFile string
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:   "yt-cli",
	Short: "A CLI for interacting with YouTube API",
	Long:  `A Dual-Mode CLI for humans and AI agents to interact with YouTube API. Supports headless authentication and deterministic output.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/yt-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output pure JSON for agents")
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	rootCmd.AddGroup(&cobra.Group{ID: "Info", Title: "Information Commands"})
	rootCmd.AddGroup(&cobra.Group{ID: "Upload", Title: "Upload Commands"})

	// Apply semantic coloring to help if interactive
	if os.Getenv("NO_COLOR") == "" && os.Getenv("YT_CLI_NO_TUI") == "" {
		usageTemplate := rootCmd.UsageTemplate()
		
		// Colorize headers
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
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		
		configDir := filepath.Join(home, ".config", "yt-cli")
		os.MkdirAll(configDir, os.ModePerm)
		
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("YT")

	if err := viper.ReadInConfig(); err == nil && !viper.GetBool("json") {
		// Outputting to stderr prevents polluting stdout JSON streams
		fmt.Fprintln(os.Stderr, ui.Muted.Render("Using config file:"), viper.ConfigFileUsed())
	}
}
