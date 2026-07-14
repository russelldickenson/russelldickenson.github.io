package config

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site     SiteConfig     `yaml:"site"`
	Build    BuildConfig    `yaml:"build"`
	Sidebar  SidebarConfig  `yaml:"sidebar"`
	Features FeaturesConfig `yaml:"features"`
	RootDir  string         `yaml:"-"`
}

type SiteConfig struct {
	BaseURL     string `yaml:"base_url"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
	Author      string `yaml:"author"`
	Favicon     string `yaml:"favicon"`
}

type BuildConfig struct {
	ContentDir   string `yaml:"content_dir"`
	OutputDir    string `yaml:"output_dir"`
	PostsPerPage int    `yaml:"posts_per_page"`
}

type SidebarConfig struct {
	Position    string `yaml:"position"`
	ShowTags    bool   `yaml:"show_tags"`
	ShowRecent  bool   `yaml:"show_recent"`
	RecentCount int    `yaml:"recent_count"`
}

type FeaturesConfig struct {
	RSS              bool `yaml:"rss"`
	CodeHighlighting bool `yaml:"code_highlighting"`
	Search           bool `yaml:"search"`
}

func Default() *Config {
	return &Config{
		Site: SiteConfig{
			Title:       "My Blog",
			Description: "A simple blog",
			URL:         "http://localhost:8080",
			Author:      "Anonymous",
		},
		Build: BuildConfig{
			ContentDir: "content/posts",
			OutputDir:  "public",
		},
		Sidebar: SidebarConfig{
			Position:    "right",
			ShowTags:    true,
			ShowRecent:  true,
			RecentCount: 5,
		},
		Features: FeaturesConfig{
			RSS:              true,
			CodeHighlighting: true,
			Search:           true,
		},
		RootDir: ".",
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()
	configPath := resolveConfigPath(path)
	cfg.RootDir = resolveRootDir(configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.applyEnvOverrides()

	return cfg, nil
}

func (c *Config) Validate() error {
	c.applyEnvOverrides()

	if c.Site.URL == "" {
		c.Site.URL = "http://localhost:8080"
	}
	c.Site.URL = strings.TrimRight(c.Site.URL, "/")

	if c.Site.BaseURL == "" {
		c.Site.BaseURL = deriveBaseURL(c.Site.URL)
	}
	c.Site.BaseURL = normalizeBaseURL(c.Site.BaseURL)

	if c.Build.ContentDir == "" {
		c.Build.ContentDir = "content/posts"
	}
	if c.Build.OutputDir == "" {
		c.Build.OutputDir = "public"
	}
	if c.Build.PostsPerPage <= 0 {
		c.Build.PostsPerPage = 10
	}
	if c.Sidebar.Position != "left" && c.Sidebar.Position != "right" {
		c.Sidebar.Position = "right"
	}
	if c.Sidebar.RecentCount < 1 {
		c.Sidebar.RecentCount = 5
	}
	return nil
}

func (c *Config) OutputPath() string {
	return c.ResolveProjectPath(c.Build.OutputDir)
}

func (c *Config) applyEnvOverrides() {
	if siteURL := strings.TrimSpace(os.Getenv("MYBLOG_SITE_URL")); siteURL != "" {
		c.Site.URL = siteURL
		if strings.TrimSpace(os.Getenv("MYBLOG_BASE_URL")) == "" {
			c.Site.BaseURL = deriveBaseURL(siteURL)
		}
	} else if siteURL := strings.TrimSpace(os.Getenv("CI_PAGES_URL")); siteURL != "" {
		c.Site.URL = siteURL
		if strings.TrimSpace(os.Getenv("MYBLOG_BASE_URL")) == "" {
			c.Site.BaseURL = deriveBaseURL(siteURL)
		}
	}

	if baseURL := strings.TrimSpace(os.Getenv("MYBLOG_BASE_URL")); baseURL != "" {
		c.Site.BaseURL = baseURL
	}
}

func (c *Config) BaseURL() string {
	return normalizeBaseURL(c.Site.BaseURL)
}

func (c *Config) RelURL(resourcePath string) string {
	base := c.BaseURL()
	trimmed := strings.TrimLeft(resourcePath, "/")
	if trimmed == "" {
		return base
	}
	if base == "/" {
		return "/" + trimmed
	}
	return base + trimmed
}

func (c *Config) AbsURL(resourcePath string) string {
	siteURL := strings.TrimRight(c.Site.URL, "/")
	if siteURL == "" {
		return c.RelURL(resourcePath)
	}

	parsed, err := url.Parse(siteURL)
	if err == nil && parsed.Path != "" && parsed.Path != "/" {
		base := siteURL + "/"
		trimmed := strings.TrimLeft(resourcePath, "/")
		if trimmed == "" {
			return base
		}
		baseURL, baseErr := url.Parse(base)
		ref, refErr := url.Parse(trimmed)
		if baseErr == nil && refErr == nil {
			return baseURL.ResolveReference(ref).String()
		}
	}

	if resourcePath == "" || resourcePath == "/" {
		return siteURL + c.BaseURL()
	}
	return siteURL + c.RelURL(resourcePath)
}

func normalizeBaseURL(raw string) string {
	if raw == "" || raw == "/" {
		return "/"
	}

	cleaned := path.Clean("/" + strings.TrimSpace(raw))
	if cleaned == "." {
		return "/"
	}
	if !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}
	return cleaned
}

func deriveBaseURL(siteURL string) string {
	parsed, err := url.Parse(siteURL)
	if err != nil || parsed.Path == "" {
		return "/"
	}
	return normalizeBaseURL(parsed.Path)
}

func (c *Config) ResolveProjectPath(target string) string {
	if target == "" || filepath.IsAbs(target) {
		return target
	}
	if c.RootDir == "" {
		return filepath.Clean(target)
	}
	return filepath.Join(c.RootDir, target)
}

func resolveRootDir(configPath string) string {
	if configPath == "" {
		return "."
	}
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return filepath.Dir(configPath)
	}
	return filepath.Dir(absPath)
}

func resolveConfigPath(configPath string) string {
	if configPath == "" || filepath.IsAbs(configPath) {
		return configPath
	}

	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	wd, err := os.Getwd()
	if err != nil {
		return configPath
	}

	for {
		candidate := filepath.Join(wd, configPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			return configPath
		}
		wd = parent
	}
}
