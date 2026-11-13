# yt-cli

A dual-mode command-line tool for humans and AI coding agents to inspect YouTube channels, list and filter channel uploads, resolve video-to-channel ownership, and upload videos safely via the YouTube Data API v3.

[![Go Reference](https://pkg.go.dev/badge/github.com/ghchinoy/yt-cli.svg)](https://pkg.go.dev/github.com/ghchinoy/yt-cli)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Authentication & Prerequisites](#authentication--prerequisites)
- [Agent Skill Installation](#agent-skill-installation)
- [Development](#development)
- [Publishing](#publishing)
- [Contributing](#contributing)
- [License](#license)

## Installation

Install the binary directly with Go (`1.26+`):

```bash
go install github.com/ghchinoy/yt-cli@latest
```

## Usage

Inspect your authenticated channel or any public channel by ID:

```bash
# View your own channel's title, ID, subscribers, and view count
yt-cli channel info --mine

# Output pure JSON for scripts or AI agents
yt-cli channel info UC_x5XG1OV2P6uZZ5FSM9Ttw --json
```

Common workflows across channels and videos:

```bash
# List the 5 most recent uploads for a channel as JSON
yt-cli channel videos UC_x5XG1OV2P6uZZ5FSM9Ttw --limit 5 --json

# List your own channel's public and unlisted videos with detailed metadata
yt-cli channel videos --mine --detailed --filter public,unlisted

# Look up which channel uploaded a specific video ID
yt-cli video channel-info dQw4w9WgXcQ --json

# Validate a video upload payload without uploading (--dry-run)
yt-cli video upload ./my-video.mp4 --title "Demo Release" --description "Walkthrough" --dry-run --json
```

| Command | Description | Key Flags |
| :--- | :--- | :--- |
| `yt-cli channel info [channel-id]` | Fetch channel title, ID, subscriber count, and total views | `--mine` (`-m`), `--json` |
| `yt-cli channel videos [channel-id]` | List uploads for a channel with optional privacy filtering | `--mine` (`-m`), `--limit` (`-l`), `--detailed` (`-d`), `--filter` (`-f`), `--json` |
| `yt-cli video channel-info <video-id>` | Resolve owning channel ID and title from a video ID (alias: `info`) | `--json` |
| `yt-cli video upload <filepath>` | Upload an `.mp4` video (defaults to `private` visibility) | `--title` (`-t`, required), `--description` (`-d`), `--dry-run`, `--json` |

Set `NO_COLOR=1` or `YT_CLI_NO_TUI=1` to disable ANSI terminal styling in non-JSON human output.

## Authentication & Prerequisites

1. Create or select a Google Cloud project and enable the [YouTube Data API v3](https://developers.google.com/youtube/v3/getting-started):
   ```bash
   gcloud services enable youtube.googleapis.com
   ```
2. Download an OAuth 2.0 **Desktop App** `client_secret.json` from the Google Cloud Console and place it in `~/.config/yt-cli/`:
   ```bash
   mkdir -p ~/.config/yt-cli
   cp /path/to/client_secret.json ~/.config/yt-cli/client_secret.json
   ```
   Alternatively, set `YT_CLIENT_ID` and `YT_CLIENT_SECRET` in your environment.
3. Run an interactive command once to complete the browser OAuth flow and cache `~/.config/yt-cli/token.json`:
   ```bash
   yt-cli channel info --mine
   ```

## Agent Skill Installation

This repository bundles an [agentskills.io](https://agentskills.io)-compliant skill at [`skills/yt-cli/SKILL.md`](skills/yt-cli/SKILL.md) alongside root agent guidance in [`GEMINI.md`](GEMINI.md).

Install the skill with [Gemini CLI](https://geminicli.com/) or the Open Agent Skills CLI:

```bash
# Gemini CLI
gemini skills install https://github.com/ghchinoy/yt-cli.git --path skills/yt-cli

# Open Agent Skills CLI
npx skills add ghchinoy/yt-cli
```

## Development

Prerequisites: **Go 1.26.11+** (and **Node.js 22.23+** if working on the documentation site in `./site`).

```bash
git clone https://github.com/ghchinoy/yt-cli.git
cd yt-cli

# Build the CLI binary
go build -o yt-cli .

# Run unit tests with the race detector
go test -p 2 -race ./...

# Run static analysis and vulnerability checks
golangci-lint run ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

To run and build the Astro landing site in [`./site`](site/):

```bash
cd site
npm install
npm run dev
ASTRO_TELEMETRY_DISABLED=1 npm run check
ASTRO_TELEMETRY_DISABLED=1 npm run build
```

## Publishing

CLI releases are distributed via Go module tags (`go install github.com/ghchinoy/yt-cli@vX.Y.Z`):

```bash
git tag -s v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

The static site in `./site` builds to `./site/dist` via `ASTRO_TELEMETRY_DISABLED=1 npm run build` for static hosting on GitHub Pages.

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change. Make sure `go test -race ./...` and `golangci-lint run ./...` pass before submitting a PR.

## License

Licensed under the [Apache License 2.0](LICENSE).

> **Disclaimer:** This is not an officially supported Google product.
