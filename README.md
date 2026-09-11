<p align="center">
  <img src="media/hero.svg" width="800" alt="simpleblog — A light/dark static blog engine written in Go" />
</p>

<p align="center">
  <b>simpleblog</b> turns Markdown posts into a complete static website with tags,
  search, RSS, and automatic light/dark themes — all from a single Go binary.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License: MIT" />
  <img src="https://img.shields.io/badge/dependencies-2%20modules-30363d" alt="Dependencies: 2 modules" />
  <img src="https://img.shields.io/badge/platform-macOS%20|%20Linux-lightgrey" alt="Platform: macOS | Linux" />
</p>

<br />

## Quick start

Write a post in `content/posts/`, then build and preview:

```bash
# Build and serve with live preview
./blog.sh -serve

# Open http://localhost:8080 in your browser
```

That's it. The engine reads your Markdown, renders it through HTML templates, and outputs a complete static site into `public/`.

To build without the dev server:

```bash
./blog.sh -build
```

<br />

## How it works

<p align="center">
  <img src="media/pipeline.svg" width="720" alt="simpleblog pipeline: Markdown posts → Go engine → static site → host anywhere" />
</p>

Posts are authored as Markdown files with YAML front matter. The Go engine parses them via [goldmark](https://github.com/yuin/goldmark), renders them through Go `html/template` layouts, and writes the full site to `public/` — ready to deploy on any static host.

<br />

## Features

<table>
  <tr>
    <td width="50%">
      <h4>☀️ / 🌙 Dual themes</h4>
      Light and dark modes that follow the OS setting, with a manual toggle. <br /><br />
      <em>See <a href="THEME.md">THEME.md</a> for the full theming guide.</em>
    </td>
    <td width="50%">
      <h4>🏷️ Tagging + tag cloud</h4>
      Tag your posts, browse via the tag cloud, and multi-select tags to filter the post list.
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h4>🔍 Client-side search</h4>
      The engine builds a search index at build time. Users can search posts without a server.
    </td>
    <td width="50%">
      <h4>📡 RSS / Atom feeds</h4>
      Auto-generated feeds so readers can subscribe. Toggle on/off in <code>config.yaml</code>.
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h4>📐 Configurable sidebar</h4>
      Position it left or right. Shows a search field, tag cloud, and recent posts list.
    </td>
    <td width="50%">
      <h4>📄 Pagination</h4>
      Configurable posts-per-page with "Newer / Older" navigation links.
    </td>
  </tr>
  <tr>
    <td width="50%">
      <h4>🤖 AI/agent discoverability</h4>
      Auto-generated <code>robots.txt</code>, <code>llms.txt</code>, and <code>sitemap.xml</code>, plus JSON-LD structured data on every post — so search engines and AI agents can index and cite your content. Toggle each on/off in <code>config.yaml</code>.
    </td>
    <td width="50%"></td>
  </tr>
</table>

<br />

## Configuration

Edit [`config.yaml`](config.yaml) to customize the blog:

| Setting | Description | Default |
|---|---|---|
| `site.title` | Blog title | `My Blog` |
| `site.description` | Tagline shown in the header | `A simple static blog built with Go` |
| `site.base_url` | Base URL for links | `/` |
| `site.author` | Author name for the footer | `Russell Dickenson` |
| `build.posts_per_page` | Posts shown per page | `4` |
| `sidebar.position` | Sidebar side (`left` or `right`) | `right` |
| `sidebar.show_tags` | Show tag cloud in sidebar | `true` |
| `sidebar.show_recent` | Show recent posts | `true` |
| `features.rss` | Enable RSS/Atom feed | `true` |
| `features.search` | Enable client-side search | `true` |
| `features.code_highlighting` | Enable syntax highlighting | `true` |
| `features.robots` | Generate `robots.txt` | `true` |
| `features.llms_txt` | Generate `llms.txt` (an LLM-readable post index) | `true` |
| `features.sitemap` | Generate `sitemap.xml` | `true` |
| `features.jsonld` | Emit JSON-LD structured data on post pages | `true` |

Flags passed on the command line take precedence over the config file:

| Flag | Description | Default |
|---|---|---|
| `-build` | Build the blog | `true` |
| `-serve` | Build and serve with a local HTTP server | — |
| `-config` | Path to config file | `config.yaml` |
| `-port` | Port for the dev server | `8080` |

### Examples

```bash
# Build only
./blog.sh -build

# Build and serve
./blog.sh -serve

# Custom port
./blog.sh -serve -port 3000

# Custom config
./blog.sh -build -config my-site.yaml
```

<br />

## Project structure

```
.
├── blog.sh                  # Platform-aware entry script (detects OS/arch)
├── blog-darwin-arm64        # Prebuilt binary (macOS, Apple Silicon)
├── blog-linux-amd64         # Prebuilt binary (Linux, x86_64)
├── config.yaml              # Site configuration
├── content/
│   └── posts/               # Markdown blog posts
├── public/                  # Generated static site output
├── src/                     # Go source code
│   ├── main.go              # CLI entry point, flag parsing, dev server
│   ├── config/              # Configuration loading & validation
│   ├── content/             # Markdown parsing (goldmark), front matter
│   ├── generator/           # HTML template rendering & site generation
│   ├── search/              # Client-side search index builder
│   ├── feed/                # RSS/Atom feed generator
│   └── discovery/           # robots.txt, llms.txt, sitemap.xml generator
├── static/                  # Static assets (JS, images, favicon)
├── themes/                  # HTML templates & CSS (light + dark)
│   ├── layout.html
│   ├── post.html
│   ├── tag.html
│   ├── sidebar.html
│   └── style.css
└── media/                   # README assets
```

<br />

## Building from source

### Prerequisites

- Go 1.26+

### Run from source

```bash
cd src
go run . -build     # Build the site
go run . -serve     # Build and serve
```

The engine finds `config.yaml` at the repo root automatically when run from `src/`.

### Build platform binaries

```bash
cd src
./build-blog-binaries.sh
```

This produces `blog-darwin-arm64` and `blog-linux-amd64` at the repo root.
The Linux binary is intended for CI/CD environments (e.g. GitLab Pages).
To build for additional platforms, uncomment the relevant lines in the build script.

<br />

## License

MIT © [Russell Dickenson](https://github.com/russelldickenson)
