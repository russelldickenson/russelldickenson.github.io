package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateNormalizesRelativeAndAbsoluteURLs(t *testing.T) {
	cfg := &Config{
		Site: SiteConfig{
			URL:     "https://example.com/myblog/",
			BaseURL: "myblog",
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if got, want := cfg.BaseURL(), "/myblog/"; got != want {
		t.Fatalf("BaseURL() = %q, want %q", got, want)
	}
	if got, want := cfg.RelURL("posts/hello.html"), "/myblog/posts/hello.html"; got != want {
		t.Fatalf("RelURL() = %q, want %q", got, want)
	}
	if got, want := cfg.AbsURL("posts/hello.html"), "https://example.com/myblog/posts/hello.html"; got != want {
		t.Fatalf("AbsURL() = %q, want %q", got, want)
	}
}

func TestValidateDerivesGitLabBaseURLFromCIPagesURL(t *testing.T) {
	t.Setenv("CI_PAGES_URL", "https://russell.gitlab.io/myblog")
	t.Setenv("MYBLOG_BASE_URL", "")

	cfg := Default()
	cfg.Site.BaseURL = "/"

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if got, want := cfg.Site.URL, "https://russell.gitlab.io/myblog"; got != want {
		t.Fatalf("Site.URL = %q, want %q", got, want)
	}
	if got, want := cfg.BaseURL(), "/myblog/"; got != want {
		t.Fatalf("BaseURL() = %q, want %q", got, want)
	}
	if got, want := cfg.AbsURL("feed.xml"), "https://russell.gitlab.io/myblog/feed.xml"; got != want {
		t.Fatalf("AbsURL() = %q, want %q", got, want)
	}
}

func TestResolveProjectPathUsesConfigDirectory(t *testing.T) {
	cfg := Default()
	cfg.RootDir = "/tmp/myblog"

	if got, want := cfg.ResolveProjectPath("content/posts"), "/tmp/myblog/content/posts"; got != want {
		t.Fatalf("ResolveProjectPath() = %q, want %q", got, want)
	}
}

func TestLoadFindsConfigInParentDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "config.yaml"), []byte("site:\n  title: Parent Config\n"), 0644); err != nil {
		t.Fatal(err)
	}
	srcDir := filepath.Join(root, "src")
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(srcDir)

	cfg, err := Load("config.yaml")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got, want := cfg.Site.Title, "Parent Config"; got != want {
		t.Fatalf("Site.Title = %q, want %q", got, want)
	}
	if got, want := cfg.RootDir, root; got != want {
		t.Fatalf("RootDir = %q, want %q", got, want)
	}
}
