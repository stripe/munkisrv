# munkisrv

A Go web server for serving [Munki](https://github.com/munki/munki) repositories with AWS CloudFront integration for secure package distribution.

## Overview

`munkisrv` is a lightweight HTTP server designed to serve Munki repositories. It provides:

- Static file serving for Munki repository catalogs, manifests, and icons
- AWS CloudFront signed URL generation for secure package downloads
- Health check endpoint for monitoring
- Graceful shutdown handling

## Features

- **Munki Repository Serving**: Serves embedded Munki repository files including catalogs, manifests, client resources, and icons
- **CloudFront Integration**: Generates signed URLs for package downloads through AWS CloudFront for enhanced security and performance
- **Health Checks**: Built-in health check endpoint at `/healthz`
- **Graceful Shutdown**: Proper signal handling for clean server shutdown
- **Configurable**: YAML-based configuration with environment variable overrides

## Architecture

The server provides three main endpoints:

- `GET /healthz` - Health check endpoint
- `GET /repo/*` - Serves static Munki repository files
- `GET /repo/pkgs/*` - Redirects to signed CloudFront URLs for package downloads

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
# Server configurations
server:
  host: "localhost"
  port: ":3000"

# CloudFront configurations
cloudfront:
  url: "https://your-distribution.cloudfront.net"
  key_id: "YOUR_CLOUDFRONT_KEY_ID"
  private_key: |
    -----BEGIN PRIVATE KEY-----
    YOUR_PRIVATE_KEY_CONTENT_HERE
    -----END PRIVATE KEY-----
```

### Environment variables

Configuration can be overridden using environment variables with the prefix `ENV_`:

- `ENV_SERVER_HOST` - Server host
- `ENV_SERVER_PORT` - Server port
- `ENV_CLOUDFRONT_URL` - CloudFront distribution URL
- `ENV_CLOUDFRONT_KEY_ID` - CloudFront key pair ID
- `ENV_CLOUDFRONT_PRIVATE_KEY` - CloudFront private key

## Installation

1. Clone the repository
2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the binary:

   ```bash
   go build -o munkisrv cmd/munkisrv/main.go
   ```

## Usage

Run the server with a configuration file:

```bash
./munkisrv -c path/to/config.yaml
```

The server will start on the configured port (default: `:3000`) and serve:

- Static repository files at `/repo/*`
- Package downloads via signed CloudFront URLs at `/repo/pkgs/*`
- Health checks at `/healthz`

## Dependencies

- `github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign` - CloudFront URL signing
- `github.com/go-chi/chi/v5` - HTTP router and middleware
- `github.com/spf13/viper` - Configuration management

## Security

The server uses AWS CloudFront signed URLs to secure package downloads. Private keys are parsed and validated at startup to ensure proper cryptographic operations.

## Health monitoring

The `/healthz` endpoint checks the availability of the embedded Munki repository by attempting to open the `catalogs/all` file. This ensures the server is properly configured and the repository is accessible.

## Munki integration

This server is designed to work with the [Munki](https://github.com/munki/munki) open-source software deployment system for macOS. It serves the repository structure that Munki clients expect:

- **Catalogs**: Software catalog definitions
- **Manifests**: Client-specific software manifests
- **Icons**: Application icons for the Munki client GUI
- **Client Resources**: Additional resources for Munki clients
- **Packages**: Software packages (served via CloudFront)

## Munki client configuration

At a minimum, configure your munki client to access the `/repo` path at your domain.

```bash
sudo defaults write /Library/Preferences/ManagedInstalls.plist SoftwareRepoURL https://<yourdomain>/repo
```
