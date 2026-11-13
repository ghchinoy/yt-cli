# YouTube CLI (yt-cli) for Agents

This CLI is designed for both human developers and autonomous AI agents.

## Core Features
- **Deterministic Output**: Every command supports a `--json` flag for machine-readable output and honors `NO_COLOR` / `YT_CLI_NO_TUI`.
- **Headless-Friendly**: Reads OAuth credentials from `~/.config/yt-cli/client_secret.json` (with cached tokens at `~/.config/yt-cli/token.json`), falling back to environment variables (`YT_CLIENT_ID`, `YT_CLIENT_SECRET`) and configuration via `~/.config/yt-cli/config.yaml`.
- **Safety First**: Destructive actions like `video upload` default to `private` visibility and support a `--dry-run` flag.
- **Agent Discovery**: Commands are grouped (`Info`, `Upload`) and documented with `Short`, `Long`, and `Example` blocks for easy parsing.

## Common Workflows for Agents

### 1. Initialize Credentials
Place your downloaded Desktop App OAuth client secret at `~/.config/yt-cli/client_secret.json`, or set the environment variables in CI/CD / headless environments:
```bash
export YT_CLIENT_ID="your-client-id"
export YT_CLIENT_SECRET="your-client-secret"
```

### 2. Get Channel Info
Retrieve metadata and subscriber/view statistics for a channel ID or the authenticated user (`--mine` and `[channel-id]` are mutually exclusive):
```bash
yt-cli channel info --mine --json
yt-cli channel info UC_x5XG1OV2P6uZZ5FSM9Ttw --json
```

### 3. List Channel Videos
List the last 5 videos from a channel (or `--mine`) with optional privacy filtering (`public,unlisted,private`) and output as JSON:
```bash
yt-cli channel videos UC_x5XG1OV2P6uZZ5FSM9Ttw --limit 5 --json
yt-cli channel videos --mine --detailed --filter public,unlisted --json
```

### 4. Upload a Video (Dry Run)
Validate parameters locally before uploading (`--title` is required):
```bash
yt-cli video upload ./my-video.mp4 --title "My Title" --dry-run --json
```

### 5. Get Channel Info from a Video
Look up the owning channel ID and title from a video ID (`info` is supported as an alias for `channel-info`):
```bash
yt-cli video channel-info dQw4w9WgXcQ --json
```

## Error Handling
Errors are printed to `stderr`. Look for `Hint:` in the error output for actionable correction steps.

## Development & Verification
Run formatting, static analysis (`.golangci.yml` v2), vulnerability scanning, and race-enabled unit tests before committing:
```bash
gofmt -s -w .
go vet ./...
golangci-lint run ./...
govulncheck ./...
go test -p 2 -race ./...
```
When building the static marketing/documentation site in `./site`:
```bash
cd site && ASTRO_TELEMETRY_DISABLED=1 npm run check && ASTRO_TELEMETRY_DISABLED=1 npm run build
```
