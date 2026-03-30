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

import "github.com/charmbracelet/lipgloss"

var (
	Accent  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Blue
	Command = lipgloss.NewStyle().Foreground(lipgloss.Color("248")) // Light Grey
	Pass    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
	Warn    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange
	Fail    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
	Muted   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dark Grey
	ID      = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // Teal/Mint
)

// FormatError provides a standard format for hints for agents
func FormatError(err error, hint string) string {
	out := Fail.Render("Error: ") + err.Error() + "\n"
	if hint != "" {
		out += Warn.Render("Hint: ") + Command.Render(hint) + "\n"
	}
	return out
}
