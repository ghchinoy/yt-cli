# Releasing and Publishing `yt-cli`

This project uses [GoReleaser](https://goreleaser.com/) via GitHub Actions to automate building, packaging, and publishing `yt-cli` binaries.

## How It Works

When you push a Git tag starting with `v` (e.g., `v0.1.0`), the GitHub Action defined in [`.github/workflows/release.yaml`](../.github/workflows/release.yaml) runs automatically:

1. Compiles static (`CGO_ENABLED=0`) Go binaries across macOS (`darwin`), Linux (`linux`), and Windows (`windows`) for `amd64` and `arm64`.
2. Uses Go linker flags (`-ldflags "-s -w -X github.com/ghchinoy/yt-cli/cmd.Version={{.Version}}"`) to inject the release tag into `yt-cli --version`.
3. Packages the binaries into `.tar.gz` (and `.zip` on Windows) archives alongside `yt-cli_<version>_checksums.txt`.
4. Generates a changelog from the commit history since the previous tag.
5. Publishes all release artifacts to GitHub Releases.

## Step-by-Step Release Guide

### 1. Ensure `main` is up to date and clean

```bash
git checkout main
git pull
git status
```

### 2. Create a Semantic Version Tag

Create an annotated or signed tag starting with `v`:

```bash
git tag -s v0.1.0 -m "Release v0.1.0"
```

### 3. Push the Tag to GitHub

Pushing the tag triggers the GoReleaser GitHub Action:

```bash
git push origin v0.1.0
```

### 4. Verify the Release

Open the repository's **Releases** or **Actions** tab on GitHub to watch the `Release` workflow complete and publish the cross-platform archives.

## Local Testing (Snapshot Build)

To test the GoReleaser build locally without publishing to GitHub:

```bash
goreleaser release --snapshot --clean
```

This builds all cross-platform binaries into the local `dist/` directory for inspection.
