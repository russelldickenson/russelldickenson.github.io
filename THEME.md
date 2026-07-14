# Theming Guide

The active color theme of the blog is toggled at runtime using JavaScript by updating a `data-theme` attribute on the root HTML element.

## How it works

### 1. Theme Files

All CSS styles and theme declarations are located in the `themes/` directory within a single stylesheet:

```plaintext
themes/
├── layout.html      # Shared HTML templates
├── post.html
├── tag.html
├── sidebar.html
└── style.css        # Unified stylesheet containing layouts and both themes
```

### 2. CSS Custom Properties

Both light and dark color schemes are defined within [style.css](file:///Users/russell/personal/repos/simpleblog/themes/style.css) using CSS custom properties (variables). This ensures that structural layouts are shared, preventing styling inconsistencies:

```css
/* themes/style.css */
:root {
    --bg-color: #ffffff;
    --text-color: #333333;
    --link-color: #0066cc;
    --header-bg: #f5f5f5;
    --border-color: #dddddd;
    --code-bg: #f4f4f4;
    --tag-bg: #e8e8e8;
    --backdrop-bg: rgba(0, 0, 0, 0.35);
    --handle-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
    --blockquote-border: var(--border-color);
    --tag-hover-color: #ffffff;
}

[data-theme="dark"] {
    --bg-color: #0d1117;
    --text-color: #c9d1d9;
    --link-color: #58a6ff;
    --header-bg: #161b22;
    --border-color: #30363d;
    --code-bg: #1f2428;
    --tag-bg: #21262d;
    --backdrop-bg: rgba(0, 0, 0, 0.5);
    --handle-shadow: 0 8px 20px rgba(0, 0, 0, 0.28);
    --blockquote-border: var(--link-color);
    --tag-hover-color: var(--bg-color);
}
```

### 3. Runtime Theme Switching

To prevent Flash of Unstyled Content (FOUC), theme loading and toggling works as follows:

1. **Synchronous Init**: An inline script in the `<head>` of [layout.html](file:///Users/russell/personal/repos/simpleblog/themes/layout.html) reads the user's preference from `localStorage` (or falls back to their operating system color scheme preference) and sets the `data-theme` attribute on `document.documentElement` before the page is rendered.
2. **Toggle Button**: Clicking the theme switcher button calls `toggleTheme()` in [theme-switcher.js](file:///Users/russell/personal/repos/simpleblog/static/theme-switcher.js), which toggles the `data-theme` attribute on `document.documentElement` and saves the selection to `localStorage`.

---

## Modifying or Adding Themes

To modify the colors of the light or dark themes, simply edit the CSS variable declarations inside [style.css](file:///Users/russell/personal/repos/simpleblog/themes/style.css):

* Edit values under `:root` to change the **light theme**.
* Edit values under `[data-theme="dark"]` to change the **dark theme**.

### Theme Variables Reference

| Variable | Purpose |
|----------|---------|
| `--bg-color` | Page background color |
| `--text-color` | Main text color |
| `--text-muted` | Muted metadata and helper text color |
| `--link-color` | Link and accent color |
| `--header-bg` | Header background |
| `--border-color` | Standard border and line color |
| `--code-bg` | Inline code and blockquote code blocks background |
| `--tag-bg` | Tag block background |
| `--backdrop-bg` | Sidebar overlay overlay background |
| `--handle-shadow` | Shadow color for sidebar grip handle |
| `--blockquote-border` | Left border accent color for blockquotes |
| `--tag-hover-color` | Text color of tag labels when hovered |

---

## Best Practices

1. **Keep variables synchronized**: If you add a new CSS custom property to light mode (`:root`), make sure to provide an override value for dark mode (`[data-theme="dark"]`).
2. **Verify contrast**: Ensure text elements meet readability guidelines against their respective theme background colors.
3. **Use semantic variables**: Avoid styling elements with hardcoded colors; map them to an existing or new theme variable in `style.css`.
