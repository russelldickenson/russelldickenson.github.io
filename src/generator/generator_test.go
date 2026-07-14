package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"myblog/config"
	"myblog/content"
)

func TestSlugify(t *testing.T) {
	if got, want := content.Slugify("Static Site Generators"), "static-site-generators"; got != want {
		t.Fatalf("Slugify() = %q, want %q", got, want)
	}
}

func TestLoadTemplatesFailsOnInvalidTemplate(t *testing.T) {
	root := t.TempDir()
	themesDir := filepath.Join(root, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatal(err)
	}

	badTemplate := filepath.Join(themesDir, "bad.html")
	if err := os.WriteFile(badTemplate, []byte(`{{define "layout"}}{{if}}{{end}}`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.Default()
	cfg.RootDir = root

	err := New(cfg).LoadTemplates()
	if err == nil {
		t.Fatal("LoadTemplates() error = nil, want invalid template error")
	}
	if !strings.Contains(err.Error(), badTemplate) {
		t.Fatalf("LoadTemplates() error = %q, want filename %q", err, badTemplate)
	}
}

func TestGenerateFailsWhenPostCannotRender(t *testing.T) {
	root := t.TempDir()
	themesDir := filepath.Join(root, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"layout.html": `{{define "layout"}}index{{end}}`,
		"style.css":    "",
	} {
		if err := os.WriteFile(filepath.Join(themesDir, name), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := config.Default()
	cfg.RootDir = root
	cfg.Build.OutputDir = "public"
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}

	gen := New(cfg)
	if err := gen.LoadTemplates(); err != nil {
		t.Fatal(err)
	}

	err := gen.Generate([]content.Post{{
		Slug:  "brittle-post",
		Title: "Brittle Post",
		Date:  time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC),
	}})
	if err == nil {
		t.Fatal("Generate() error = nil, want missing post template error")
	}
	if !strings.Contains(err.Error(), "generate post brittle-post") {
		t.Fatalf("Generate() error = %q, want post slug context", err)
	}
}
