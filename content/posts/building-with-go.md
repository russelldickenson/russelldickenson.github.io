---
title: "Building a Blog with Go"
date: 2024-02-10
tags: ["go", "tutorial", "web-development"]
description: "How I built this blog engine"
---

# Building a Blog with Go

I built this blog engine from scratch using Go. Here's how it works:

## Architecture

```
myblog/
├── cmd/           # CLI entry point
├── internal/       # Core packages
│   ├── config/    # Configuration
│   ├── content/   # Markdown parsing
│   ├── generator/ # HTML generation
│   ├── search/   # Search index
│   └── feed/     # RSS/Atom
├── themes/         # Theme templates
└── content/      # Blog posts
```

## Key Packages

### Markdown Parsing

Using `goldmark` for Markdown to HTML conversion with extensions for:
- GFM (GitHub Flavored Markdown)
- Syntax highlighting
- Tables and task lists

### Templating

Go's built-in `html/template` package provides safe templating with:
- Layouts and partials
- Custom functions
- Automatic escaping

### Frontmatter

YAML frontmatter in each post provides metadata:
- Title, date, tags
- Description
- Draft status

## Results

The blog generates static HTML that's fast, secure, and easy to host anywhere!