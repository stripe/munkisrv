# munkisrv

[munkisrv](https://pkg.go.dev/github.com/stripe/munkisrv) is a Go web server for serving [Munki](https://github.com/munki/munki) repositories. It comes with built-in support for generating AWS CloudFront pre-signed URLs.

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
2. Build the binary:

   ```bash
   go build ./cmd/munkisrv
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

Send a test request:

```bash
curl http://127.0.0.1:3000/repo/catalogs/all
```

Configure munki to connect:

```bash
sudo defaults write /Library/Preferences/ManagedInstalls.plist SoftwareRepoURL http://127.0.0.1:3000/repo
```

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

Ensure munki is configured to follow HTTP redirects.

```bash
sudo defaults write /Library/Preferences/ManagedInstalls.plist FollowHTTPRedirects https
```
