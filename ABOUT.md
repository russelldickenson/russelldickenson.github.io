simpleblog (module myblog)

  A custom static-site generator for a personal blog, written in Go by Russell Dickenson. It turns Markdown posts into a complete static HTML website
  that can be built once or built-and-served locally.

  What it does

  You write posts as Markdown files in content/posts/, run the engine, and it generates a finished static site into public/ — ready to host anywhere
  (the .gitlab-ci.yml suggests GitLab Pages CI/CD deployment).

  How it's organized

  The Go source (src/) is split into focused packages:

  ┌────────────┬──────────────────────────────────────────────────────────────────────────────────────────────────────────┐
  │  Package   │                                              Responsibility                                              │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ main.go    │ CLI entry — flag parsing (-build, -serve, -config, -port), orchestration, and a small static file server │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ config/    │ Loads/validates config.yaml                                                                              │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ content/   │ Parses Markdown posts (via goldmark), front matter, slugs, tags                                          │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ generator/ │ Renders posts through HTML templates into the output site                                                │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ search/    │ Builds a client-side search index                                                                        │
  ├────────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ feed/      │ Generates RSS feed                                                                                       │
  └────────────┴──────────────────────────────────────────────────────────────────────────────────────────────────────────┘
  
  Supporting assets live outside src/: themes/ (HTML templates + light/dark CSS), static/ (JS for search, theme switching, sidebar toggle, RSS), and
  content/ (the Markdown posts).

  Features (from README + config)
  
  - Reverse-chronological post listing with pagination (posts_per_page)
  - Tagging + a tag cloud; click a tag to filter posts
  - Configurable sidebar (left/right) with tag cloud, search box, and recent-posts list
  - Light/dark themes that follow the OS setting but can be toggled manually
  - Client-side search, RSS feed, and code syntax highlighting (all toggleable in config.yaml)

  Distribution

  Prebuilt binaries for macOS (Apple Silicon) and Linux (amd64, for CI) are committed at the repo root, with blog.sh as a wrapper that picks the right
  binary for the current platform. Built with Go 1.26, depending only on goldmark (Markdown) and yaml.v3.

  In short: a lightweight, self-contained, dependency-minimal blog engine — comparable to Hugo or Jekyll but hand-rolled and tailored to the author's
  own needs.
  