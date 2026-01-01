# bluenode-helper

Backend crafted for BlueNode NAS OS

## Features

- HTTP API over Unix domain sockets
- Graceful shutdown support
- Systemd integration (when installed via RPM)

## Installation

### From Binary

1. Download the latest release binary from [GitHub Releases](https://github.com/BlueNode-NAS/BNHelper/releases)
2. Make it executable: `chmod +x bluenode-helper`
3. Move to a directory in your PATH: `sudo mv bluenode-helper /usr/local/bin/`

### From RPM (Fedora/RHEL/CentOS)

```bash
sudo rpm -ivh bluenode-helper-*.rpm
sudo systemctl enable --now bluenode-helper
```

## Building from Source

### Prerequisites

- Go 1.25.4 or later

### Build

```bash
go build -o bluenode-helper
```

### Build with Version Information

```bash
go build -ldflags "-X 'main.Version=1.0.0' -X 'main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)' -X 'main.GitCommit=$(git rev-parse --short HEAD)'" -o bluenode-helper
```

## Usage

### Start the service

```bash
./bluenode-helper
```

### Check version

```bash
./bluenode-helper -version
```

## Development

### GitHub Actions

This project uses GitHub Actions for continuous integration and release automation:

- **Build Job**: Builds the binary and RPM package on every push/PR
- **Release Job**: Automatically creates a GitHub release when a commit to `main` contains `(RELEASE)` in the message

### Creating a Release

1. Update the version in `version.go`
2. Commit your changes
3. Create a commit with `(RELEASE)` in the message:
   ```bash
   git commit -m "Release version 1.0.0 (RELEASE)"
   git push origin main
   ```
4. GitHub Actions will automatically:
   - Build the binary and RPM
   - Create a git tag (e.g., `v1.0.0`)
   - Create a GitHub release with artifacts

## License

Proprietary - BlueNode NAS OS

