# Munki client configuration guide

This document explains how to configure your local Munki client to add new software packages to the munkisrv repository. This setup allows you to import packages and create pkginfo files that will be served by the munkisrv server.

## Overview

The munkisrv project uses a hybrid approach where:

- **Repository metadata** (catalogs, manifests, pkginfo files) is embedded in the server binary
- **Package files** (`.dmg`, `.pkg`, etc.) are served through CloudFront with signed URLs
- **Local development** requires configuring munkiimport to work with the local repository structure

## Prerequisites

Before configuring your Munki client, ensure you have:

1. **Munki tools installed** - Download from [Munki releases](https://github.com/munki/munki/releases/latest)
2. **Access to the munkisrv repository** - You need the full path to the `munkirepo` directory

## Repository structure

The munkisrv project uses this repository structure:

```
munkisrv/
├── munkirepo/           # This is your Munki repository
│   ├── catalogs/        # Munki catalogs
│   ├── manifests/       # Munki manifests
│   ├── pkgsinfo/        # Package information files
│   ├── icons/           # Application icons
│   └── client_resources/ # Client-side resources
└── cmd/munkisrv/        # Server application
```

## Configuration steps

### Step 1: Determine your repository path

First, identify the full path to your `munkirepo` directory. For example:

- If this project is at `/Users/foo/munkisrv/`, then your repo path is `/Users/foo/munkisrv/munkirepo`
- If this project is at `/opt/munkisrv/`, then your repo path is `/opt/munkisrv/munkirepo`

### Step 2: Configure munkiimport preferences

Run these commands in Terminal, replacing `$REPO_PATH` with your actual repository path:

```bash
# Set the repository plugin type
defaults write com.googlecode.munki.munkiimport plugin FileRepo

# Set the repository URL (use file:// protocol for local development)
defaults write com.googlecode.munki.munkiimport repo_url "file://$REPO_PATH"

# Set the default catalog for new packages
defaults write com.googlecode.munki.munkiimport default_catalog testing

# Set the pkginfo file extension
defaults write com.googlecode.munki.munkiimport pkginfo_extension .plist
```

### Example configuration

For a typical setup with the repository at `/Users/username/munkisrv/munkirepo`:

```bash
defaults write com.googlecode.munki.munkiimport plugin FileRepo
defaults write com.googlecode.munki.munkiimport repo_url "file:///Users/username/munkisrv/munkirepo"
defaults write com.googlecode.munki.munkiimport default_catalog testing
defaults write com.googlecode.munki.munkiimport pkginfo_extension .plist
```

### Step 3: Verify configuration

You can verify your settings with:

```bash
# Check all munkiimport preferences
defaults read com.googlecode.munki.munkiimport

# Or check individual settings
defaults read com.googlecode.munki.munkiimport repo_url
defaults read com.googlecode.munki.munkiimport default_catalog
```

## Adding new software

### Step 1: Prepare your package

1. **Download or create your package** (`.dmg`, `.pkg`, etc.)
2. **Place it in a temporary location** for import

### Step 2: Import the package

Use the `munkiimport` command to add your package:

```bash
# Basic import
/usr/local/munki/munkiimport /path/to/your/package.dmg

# Import with specific options
/usr/local/munki/munkiimport \
  --catalog testing \
  --name "YourAppName" \
  --displayname "Your App Display Name" \
  /path/to/your/package.dmg
```

### Step 3: Edit the pkginfo file

After import, edit the generated pkginfo file to ensure proper configuration:

```bash
# Find the generated pkginfo file
find /path/to/munkirepo/pkgsinfo -name "*YourAppName*" -type f

# Edit the file (replace with your preferred editor)
open /path/to/munkirepo/pkgsinfo/apps/yourapp/yourapp-1.0.plist
```

### Step 4: Update the repository

After making changes, update the repository:

```bash
# Make catalogs from pkginfo files
/usr/local/munki/makecatalogs -s /path/to/munkirepo

# Or if you're in the repository directory
cd /path/to/munkirepo
/usr/local/munki/makecatalogs
```

## Important considerations

### Package file handling

Since munkisrv serves packages through CloudFront, you have two options:

1. **Local Development**: Keep packages in the local `munkirepo/pkgs/` directory for testing
2. **Production**: Upload packages to S3

### `installer_item_location` configuration

For packages served through CloudFront, ensure the `installer_item_location` in your pkginfo file matches your S3 upload path:

```xml
<key>installer_item_location</key>
<string>apps/yourapp/yourapp-1.0.dmg</string>
```

This path should correspond to your S3 upload path: `s3://your-bucket/<your-s3-prefix>/apps/yourapp/yourapp-1.0.dmg`

### Catalog management

- **testing catalog**: Use for new packages during development
- **all catalog**: Use for production-ready packages
- **Update catalogs** Run `makecatalogs` after making changes to pkginfo files
