# AGENTS.md

Guidance for AI coding agents working in this repo. Read this before making changes.

## What this is

`simpleblog` is a static-site generator written in Go: Markdown posts in →
templated static HTML site out (`public/`). No database, no server-side
rendering at request time, no build step beyond the Go binary itself. Keep
that shape — don't introduce a backend, database, or SSR unless explicitly
asked.

## Source layout

```
src/main.go              CLI entry point, flag parsing, dev server (-serve)
src/config/               config.yaml loading, validation, URL helpers (RelURL/AbsURL)
src/content/               Markdown + YAML frontmatter parsing (goldmark)
src/generator/             renders themes/*.html to public/, pagination, tag pages
src/search/                builds public/search.json for client-side search
src/feed/                  RSS (feed.xml) + Atom (feed.atom) generation
themes/*.html               Go html/template files (layout, post, tag, 404, sidebar)
themes/style.css            single stylesheet, CSS custom properties for light/dark
static/*.js                 vanilla JS, no framework, no bundler, no build step
content/posts/*.md          blog posts (YAML frontmatter + Markdown body)
config.yaml                 site config (see README.md for the full key reference)
```

`go.mod` module name is `myblog`; run Go commands from `src/`.

## Build, test, run

```bash
cd src
go build ./...          # compile
go vet ./...             # static checks
go test ./...             # unit tests (config, content, generator, search)
go run . -build           # generate public/ from content/posts/
go run . -serve -port 8080  # generate + serve public/ locally
```

Or via the wrapper script from the repo root: `./blog.sh -serve`. The
prebuilt `blog-darwin-arm64` / `blog-linux-amd64` binaries are stale as soon
as source changes — rebuild before trusting them (`cd src && go build -o
../blog-darwin-arm64 .` or similar, matching your platform).

`go run .` / the binary looks for `config.yaml` relative to its working
directory by walking up, so running from `src/` still finds the repo-root
config.

There is no JS test suite and no linter configured for `static/*.js` —
verify JS changes by building the site and checking behavior in a browser
(see "Verifying changes" below).

## Conventions and gotchas specific to this repo

- **Templates auto-escape by default** (`html/template`). Only post body HTML
  (`post.HTML`, already sanitized by goldmark with `Unsafe` *not* set) is
  passed through as `template.HTML`. Never mark frontmatter fields (title,
  description, tags) as `template.HTML` — they must stay auto-escaped.
- **Client-side rendering must never use `innerHTML` with unescaped
  interpolation.** `static/tags.js` had a reflected XSS via `?tags=` fixed by
  switching to `document.createElement`/`textContent`/`setAttribute` instead
  of template-literal `innerHTML`. Follow that pattern for any new dynamic
  DOM rendering — build nodes, don't concatenate HTML strings from
  data that isn't 100% static.
- **`search.json` is intentionally minimal** — slug/title/date/tags/description
  only. Don't add post body content back into it; neither `search.js` nor
  `tags.js` need it, and it would just bloat the payload every visitor
  downloads. If you need full-text search over post bodies, that's a
  deliberate feature decision to discuss, not a silent addition.
- **`search.js` and `tags.js` share one fetch of `search.json`** via
  `static/search-index.js` (`window.blogSearchIndex.load(baseURL)`, cached
  promise). Any new script that needs the search index should use this
  loader rather than fetching `search.json` directly — don't reintroduce a
  duplicate fetch.
- **Pagination**: the index page and tag pages both paginate using
  `build.posts_per_page` from config, with matching `PrevLink`/`NextLink`/
  `CurrentPage`/`TotalPages` fields on `generator.PageData`. If you add a new
  listing page (e.g. an author page), mirror the existing pattern in
  `generateIndex`/`generateTagPages` in `src/generator/generator.go` rather
  than inventing a new scheme.
- **`RelURL` vs `AbsURL`** (`src/config/config.go`): `RelURL` rewrites a
  path under `base_url` and passes through absolute URLs (`http(s)://`,
  `//`, `mailto:`, etc.) unchanged — use it for any in-template link that
  might reasonably be either internal or external (e.g. `header_links`).
  `AbsURL` always produces a fully-qualified URL against `site.url` — use it
  for anything that leaves the site's own pages (RSS/Atom, canonical links,
  OG tags).
- **RSS/Atom date formats are not interchangeable.** RSS 2.0 requires
  RFC822-style dates (`time.RFC1123` in Go); Atom requires RFC3339. Don't
  copy a date field from one feed's format into the other.
- **The dev server (`-serve`) is not a general-purpose file server** — it
  exists for local preview. It cleans the URL path before any filesystem
  stat (`path.Clean`) to avoid leaking file existence outside `public/`, and
  it has request timeouts set. If you touch `serve()` in `src/main.go`,
  preserve both properties; don't reintroduce a bare `http.Server{}`/
  `http.ListenAndServe` or an unclean path join.
- **404 handling**: `themes/404.html` is a real template rendered to
  `public/404.html` at build time (via `generator.generate404`), using the
  same header/sidebar/footer chrome as other pages. The dev server serves it
  with a genuine `404` status for both missing `.html` paths and
  extensionless paths (see `notFoundWriter` in `src/main.go`). Static hosts
  (GitHub Pages, Netlify, etc.) pick up `public/404.html` by convention —
  don't rely on the dev server's behavior alone.
- **Draft posts** (`draft: true` in frontmatter) are parsed but excluded from
  every generated page, the search index, and both feeds. Preserve that
  filter in any new generation path that iterates posts.
- **Client-side tag filtering vs server-rendered tag pages**: `tags.js`
  reads/writes the `?tags=` query param and re-renders the post list
  client-side (works across the index and tag pages without a reload); the
  server-rendered `themes/tag.html` output is the no-JS/crawler fallback and
  must stay correct (including pagination) independent of the JS layer.
  Changes to one should be checked against the other.

## Verifying changes

This is a visual, static-HTML product — passing `go test` is necessary but
not sufficient for anything touching templates, CSS, or client JS.

1. `cd src && go build ./... && go vet ./... && go test ./...`
2. `go run . -build` (or `-serve`) and inspect the generated output in
   `public/` directly (`curl`/`cat` the relevant HTML/XML/JSON file) for
   anything server-rendered — templates, feeds, search index, pagination.
3. For anything interactive (search, tag filtering, theme toggle, sidebar,
   404 page) or visual (CSS, responsive layout, light/dark mode), actually
   load the page in a browser and exercise it — don't infer correctness from
   reading the JS alone. Check both light and dark mode
   (`prefers-color-scheme` / the theme toggle) since `themes/style.css` uses
   CSS custom properties that differ per theme.
4. If a change is security-sensitive (anything rendering user- or
   content-derived strings into HTML, the dev server's path handling), prove
   it with an actual payload/request, not just code review — e.g. navigate
   to a crafted URL and confirm no script executes, or `curl` a traversal
   attempt and confirm it's rejected.
5. Clean up: kill any dev server you started, don't commit `public/` (it's
   gitignored — regenerated output, not source).

## Git / commit conventions

- Small, focused commits per fix or feature; imperative mood subject line
  (see `git log` for the house style), with a body explaining *why* when the
  fix isn't self-evident from the diff.
- Don't commit generated output (`public/`) or scratch/debug files.
- Don't rebuild and commit the prebuilt platform binaries
  (`blog-darwin-arm64`, `blog-linux-amd64`) as part of unrelated source
  changes — that's a separate, deliberate release step
  (`src/build-blog-binaries.sh`).
