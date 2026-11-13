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

// Package ui provides semantic terminal styling and error formatting for yt-cli.
package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Accent highlights navigation landmarks, headers, and group titles.
	Accent = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	// Command styles scan targets such as command names, flags, and file paths.
	Command = lipgloss.NewStyle().Foreground(lipgloss.Color("248"))
	// Pass indicates success states and completed tasks.
	Pass = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	// Warn indicates transient states, dry-run notices, and actionable hints.
	Warn = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	// Fail indicates errors and rejected states.
	Fail = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	// Muted de-emphasizes secondary metadata and defaults.
	Muted = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// ID highlights unique resource identifiers.
	ID = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
)

// FormatError formats an error and optional actionable hint for stderr output.
func FormatError(err error, hint string) string {
	if err == nil {
		return ""
	}
	out := Fail.Render("Error: ") + err.Error() + "\n"
	if hint != "" {
		out += Warn.Render("Hint: ") + Command.Render(hint) + "\n"
	}
	return out
}
