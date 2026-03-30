# YouTube CLI (yt-cli) for Agents

This CLI is designed for both human developers and autonomous AI agents.

## Core Features
- **Deterministic Output**: Every command supports a `--json` flag for machine-readable output.
- **Headless-Friendly**: Supports environment variables (`YT_CLIENT_ID`, `YT_CLIENT_SECRET`) and configuration via `~/.config/yt-cli/config.yaml`.
- **Safety First**: Destructive actions like `video upload` support a `--dry-run` flag.
- **Agent Discovery**: Commands are structured and documented with examples for easy parsing.

## Common Workflows for Agents

### 1. Initialize Credentials
If running in a CI/CD environment or as an agent, set the environment variables:
```bash
export YT_CLIENT_ID="your-client-id"
export YT_CLIENT_SECRET="your-client-secret"
```

### 2. List Channel Videos
To list the last 5 videos from a channel and output as JSON:
```bash
yt-cli channel videos UC... --limit 5 --json
```

### 3. Upload a Video (Dry Run)
Validate parameters before uploading:
```bash
yt-cli video upload ./my-video.mp4 --title "My Title" --dry-run
```

### 4. Get Channel info from a Video
```bash
yt-cli video info dQw4w9WgXcQ --json
```

## Error Handling
Errors are printed to `stderr`. Look for `Hint:` in the error output for actionable correction steps.
