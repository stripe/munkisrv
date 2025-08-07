# S3 upload and CloudFront configuration guide

This document explains how to configure S3 uploads and CloudFront for the munkisrv project. The system is designed to serve Munki packages through CloudFront with signed URLs.

## Overview

The munkisrv project serves Munki repository metadata (catalogs, manifests, etc.) directly from embedded files, but package files (`.dmg`, `.pkg`, etc.) are served through CloudFront with signed URLs.

## S3 upload process

### Upload path requirements

**Critical**: Your S3 upload path must match the `installer_item_location` value in your pkginfo files.

For example, if your pkginfo file contains:

```xml
<key>installer_item_location</key>
<string>apps/google/chrome/googlechrome-133.0.6943.127.dmg</string>
```

Then your S3 upload path should be:

```bash
s3://your-bucket-name/repo/pkgs/apps/google/chrome/googlechrome-133.0.6943.127.dmg
```

### Upload process

The upload process varies by organization, but here are common approaches:

#### Option 1: AWS CLI
```bash
# Upload a package to S3
aws s3 cp googlechrome-133.0.6943.127.dmg s3://your-bucket-name/repo/pkgs/apps/google/chrome/

# Upload with specific path matching pkginfo
aws s3 cp googlechrome-133.0.6943.127.dmg s3://your-bucket-name/repo/pkgs/apps/google/chrome/googlechrome-133.0.6943.127.dmg
```

#### Option 2: AWS SDK/API

Use your preferred AWS SDK to upload files to the correct S3 path.

#### Option 3: CI/CD pipeline

Integrate S3 uploads into your CI/CD pipeline to automatically upload packages when they're built.

### S3 Bucket configuration

1. **Create an S3 bucket** for your Munki packages
2. **Configure bucket permissions** to allow CloudFront access
3. **Set up CORS** if needed for web-based uploads
4. **Consider lifecycle policies** for cost management

## CloudFront configuration

### Distribution setup

1. **Create a CloudFront distribution** pointing to your S3 bucket
2. **Configure the origin path** to match your S3 prefix (e.g., `/repo/pkgs/`)
3. **Set up behaviors** to handle package file requests

### Origin configuration

Configure your CloudFront distribution with:

- **Origin Domain**: Your S3 bucket
- **Origin Path**: `repo/pkgs` (or your chosen prefix after /repo/)
- **Origin Access**: Use Origin Access Control (OAC) or Origin Access Identity (OAI)

## Munkisrv configuration

### CloudFront settings

Update your `config.yaml` with your CloudFront settings:

```yaml
cloudfront:
  url: "https://your-distribution-id.cloudfront.net"
  key_id: "YOUR_CLOUDFRONT_KEY_ID"
  private_key: |
    -----BEGIN PRIVATE KEY-----
    Your CloudFront private key content
    -----END PRIVATE KEY-----
```

### CloudFront key pair setup

1. **Create a CloudFront key pair** in the AWS Console
2. **Download the private key** and add it to your configuration
3. **Note the Key ID** and add it to your configuration
4. **Configure the key pair** in your CloudFront distribution

## How It Works

### Request Flow

1. **Munki client** requests a package: `/repo/pkgs/apps/google/chrome/googlechrome-133.0.6943.127.dmg`
2. **munkisrv** receives the request and constructs the full CloudFront URL
3. **munkisrv** generates a signed URL** with a 1-hour expiration
4. **munkisrv** redirects the client to the signed CloudFront URL
5. **CloudFront** serves the package file directly to the client

### Code Implementation

The `munkiPkgFunc` in `cmd/munkisrv/main.go` handles this process:

```go
func munkiPkgFunc(cloudFrontURL string, signer *sign.URLSigner) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse the CloudFront base URL
        u, err := url.Parse(cloudFrontURL)
        if err != nil {
            http.Error(w, "failed to parse base url", http.StatusInternalServerError)
            return
        }
        
        // Construct the full path
        u.Path = path.Join(u.Path, r.URL.Path)
        finalURL := u.String()

        // Generate signed URL with 1-hour expiration
        signedURL, err := signer.Sign(finalURL, time.Now().Add(1*time.Hour))
        if err != nil {
            http.Error(w, "Failed to sign url", http.StatusInternalServerError)
            return
        }
        
        // Redirect to signed URL
        http.Redirect(w, r, signedURL, http.StatusTemporaryRedirect)
    }
}
```
