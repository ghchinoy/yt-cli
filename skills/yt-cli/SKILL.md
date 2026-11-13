---
name: yt-cli
description: Interact with the YouTube Data API v3 using yt-cli to inspect channels, list and filter channel videos, resolve video-to-channel ownership, and upload videos safely with --dry-run and --json output. Use when the user asks to query YouTube channels, list YouTube uploads, inspect video metadata, or upload a video from the command line.
license: Apache-2.0
metadata:
  version: "1.0.0"
---

# YouTube CLI (`yt-cli`) Skill

`yt-cli` is a dual-mode command-line tool for querying and managing YouTube channels and videos via the YouTube Data API v3. Always pass `--json` when invoking `yt-cli` as an agent so `stdout` returns structured JSON while diagnostics and hints go to `stderr`.

## Prerequisites & Authentication

1. Ensure `yt-cli` is available in your `PATH` or current working directory (`go install github.com/ghchinoy/yt-cli@latest`, or run the bundled `scripts/install.sh` to download and verify the latest prebuilt release binary).
2. Credentials are resolved in order:
   - OAuth 2.0 Desktop client secret at `~/.config/yt-cli/client_secret.json` with cached token at `~/.config/yt-cli/token.json`.
   - Environment variables `YT_CLIENT_ID` and `YT_CLIENT_SECRET`.
3. If an interactive browser flow is required and no cached `token.json` exists, ask the user to run `yt-cli channel info --mine` once in their terminal to complete OAuth consent.

## Core Commands

### 1. Channel Metadata (`channel info`)
Retrieve channel title, ID, subscriber count, and total view count. Pass **either** a positional `<channel-id>` **or** `--mine` (never both):

```bash
# Authenticated user's own channel
yt-cli channel info --mine --json

# Specific channel by ID
yt-cli channel info UC_x5XG1OV2P6uZZ5FSM9Ttw --json
```

### 2. List Channel Uploads (`channel videos`)
Fetch recent uploads from a channel's uploads playlist with `--limit` (default `10`), `--detailed`, and optional comma-separated `--filter` (`public`, `unlisted`, `private`):

```bash
# Latest 5 videos for a channel as JSON
yt-cli channel videos UC_x5XG1OV2P6uZZ5FSM9Ttw --limit 5 --json

# Authenticated owner's public and unlisted videos
yt-cli channel videos --mine --detailed --filter public,unlisted --json
```

### 3. Resolve Channel from Video (`video channel-info`)
Look up the channel ID and channel title that uploaded a given `<video-id>` (`video info` is also accepted as an alias):

```bash
yt-cli video channel-info dQw4w9WgXcQ --json
```

### 4. Upload a Video (`video upload`)
Uploads default to `private` visibility (`CategoryId: 22`). **Always** validate parameters with `--dry-run` before performing a live upload:

```bash
# Step 1: Validate payload with --dry-run
yt-cli video upload ./demo.mp4 --title "Release Demo" --description "Walkthrough" --dry-run --json

# Step 2: Execute live upload after user confirmation
yt-cli video upload ./demo.mp4 --title "Release Demo" --description "Walkthrough" --json
```

## Error Handling
When a command exits non-zero, inspect `stderr` for the `Error:` and `Hint:` lines and follow the suggested remediation before retrying.
