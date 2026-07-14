# myblog

Custom blog engine written in Go.

## Features

- Main page lists all blog posts in reverse chronological order.
- Posts can be tagged.
- Sidebar (can be positioned on the left or right) - includes tag cloud, search field, list of 5 most recent blog posts.
- Tag cloud lists all tags. Selecting a tag lists all blog posts with that tag.
- Light and dark themes. Selected automatically according to the system setting, but also can be toggled by the user.

## Quick Start

### Run the blog (serves on http://localhost:8080)

```bash
./blog.sh -serve
```

### Build only (outputs to public/)

```bash
./blog.sh -build
```

## Usage

```bash
./blog.sh [options]
```

### Options

Flags take precedence over options specified in the configuration file.

| Flag | Description | Default |
|------|-------------|---------|
| `-build` | Build the blog | `true` |
| `-serve` | Build and serve the blog | `false` |
| `-config` | Path to config file | `config.yaml` |
| `-port` | Port to serve on | `8080` |

## Configuration

Edit `config.yaml` to customize:

- Site title, description, URL
- Posts per page
- Sidebar position and content
- Theme selection (light/dark)

### Examples

```bash
# Build the blog
./blog.sh -build

# Build and serve locally
./blog.sh -serve

# Serve on custom port
./blog.sh -serve -port 3000
```

## Directory Structure

| Directory | Contents |
|-----------|----------|
| repo root | Compiled binaries for supported platforms and `blog.sh` |
| `content/` | Blog posts (Markdown files) |
| `public/` | Generated blog output |
| `src/` | Go source code and Go module files |
| `static/` | Static assets (JS, images) |
| `themes/` | HTML templates and CSS themes |

## Building from Source

### Prerequisites

- Go 1.26+

### Run from source

```bash
cd src
go run . -build
go run . -serve
```

When run from `src/`, the engine finds the root `config.yaml` automatically.

### Build binaries

The `src/build-blog-binaries.sh` script builds binaries for macOS (Silicon) and Linux at the repo
root. The Linux build is intended for CI/CD environments. To build for other platforms, uncomment
the relevant lines in the build
script.

```bash
# Build all platform binaries
cd src
./build-blog-binaries.sh
```
