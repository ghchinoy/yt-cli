# YouTube CLI (yt-cli)

An Agent-Aware Command-Line Interface for interacting with the YouTube API.


## Core Features
- **Dual-Mode Experience**: Built for both human operators (rich TUI) and AI agents (deterministic `--json` output).
- **Headless Auth**: Designed to gracefully handle authentication in autonomous environments.
- **Mutative Safety**: Built-in `--dry-run` flags for destructive actions.
- **Agent Skill**: Includes an `agentskills.io` compliant `SKILL.md` for seamless discovery by other agents.

## Installation

You can install `yt-cli` directly using Go:

```bash
go install github.com/ghchinoy/yt-cli@latest
```

## Getting Started

1. Download your OAuth 2.0 `client_secret.json` from the Google Cloud Console (Desktop App type).
2. Place it in the XDG configuration directory:
   ```bash
   mkdir -p ~/.config/yt-cli
   cp path/to/your/client_secret.json ~/.config/yt-cli/client_secret.json
   ```
3. Run any command to start the initial interactive OAuth flow:
   ```bash
   yt-cli channel info --mine
   ```

## Agent Usage

This repository includes an Agent Skill that provides instructions on how to use `yt-cli`.

If you are using [Gemini CLI](https://geminicli.com/), you can install the skill directly from this repository:

```bash
gemini skills install https://github.com/ghchinoy/yt-cli.git --path skills/yt-cli
```

Once installed, your agent will understand how to list your videos, look up channel information, and securely upload new content! See the `skills/yt-cli/SKILL.md` file for full instructions.

## Prerequisites

* Google Cloud Project
* [YouTube Data API v3](https://developers.google.com/youtube/v3/getting-started) enabled (`gcloud enable services youtube.googleapis.com`)
* YouTube Key and/or OAuth Client ID (preferred)
* A video to upload (optional)


## Disclaimer
This is not an official Google project.


Made with ❤️ by Gemini CLI and agentskills.