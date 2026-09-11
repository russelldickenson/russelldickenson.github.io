# Changelog

## [Unreleased]

### Optimized
- **Template Loading:** Refactored `LoadTemplates` in [**`src/generator/generator.go`**](file:///Users/russell/personal/repos/myblog/src/generator/generator.go#L89-L115) to parse template files directly into the template set, removing the redundant `Clone()` calls in the file parsing loop to reduce CPU overhead and memory allocations.
- **HTML Escaper:** Replaced the custom `escapeHTML` function in [**`src/generator/generator.go`**](file:///Users/russell/personal/repos/myblog/src/generator/generator.go) with the Go standard library's highly optimized and secure [**`html.EscapeString`**](https://pkg.go.dev/html#EscapeString).
- **Slugify Unification:** Consolidated duplicate `slugify` helper functions by exporting `Slugify` from [**`src/content/post.go`**](file:///Users/russell/personal/repos/myblog/src/content/post.go) and refactoring [**`src/generator/generator.go`**](file:///Users/russell/personal/repos/myblog/src/generator/generator.go) (and their tests) to call it directly.
- **CSS De-duplication:** Removed redundant and duplicate CSS rules for `.site-header h1`, `.site-header p`, and `.site-description` in [**`themes/light.css`**](file:///Users/russell/personal/repos/myblog/themes/light.css) and [**`themes/dark.css`**](file:///Users/russell/personal/repos/myblog/themes/dark.css).
- **Theme Variables Alignment:** Standardized variable names in [**`themes/dark.css`**](file:///Users/russell/personal/repos/myblog/themes/dark.css) to match properties in [**`themes/light.css`**](file:///Users/russell/personal/repos/myblog/themes/light.css) and guidelines in [**`THEME.md`**](file:///Users/russell/personal/repos/myblog/THEME.md), and defined the missing `--text-muted` property in the light theme variables block.

### Removed
- **Unused Helper:** Removed the unused `copyDir` helper function from [**`src/generator/generator.go`**](file:///Users/russell/personal/repos/myblog/src/generator/generator.go).