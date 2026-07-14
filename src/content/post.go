package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	Draft       bool      `json:"draft"`
	Content     string    `json:"content"`
	HTML        string    `json:"html"`
}

type ByDate []Post

func (p ByDate) Len() int           { return len(p) }
func (p ByDate) Less(i, j int) bool { return p[i].Date.After(p[j].Date) }
func (p ByDate) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Parser struct {
	md goldmark.Markdown
}

func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
			extension.TaskList,
			extension.Strikethrough,
		),
	)
	return &Parser{md: md}
}

func (p *Parser) ParseFile(path string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return p.Parse(string(data), path)
}

func (p *Parser) Parse(content, path string) (*Post, error) {
	var fmData map[string]interface{}
	var body string

	lines := strings.SplitN(content, "\n", 2)
	if len(lines) < 2 {
		return &Post{
			Content: content,
			Slug:    Slugify(filepath.Base(path)),
		}, nil
	}

	if strings.TrimSpace(lines[0]) == "---" && strings.Contains(lines[1], "---") {
		fmAndBody := strings.SplitN(content, "---", 3)
		if len(fmAndBody) >= 3 {
			fmPart := fmAndBody[1]
			body = strings.TrimPrefix(fmAndBody[2], "\n")

			if err := yaml.Unmarshal([]byte(fmPart), &fmData); err != nil {
				return nil, fmt.Errorf("frontmatter parse error: %w", err)
			}
		}
	}

	post := &Post{
		Slug: Slugify(filepath.Base(path)),
	}

	if title, ok := fmData["title"].(string); ok {
		post.Title = title
	}
	if desc, ok := fmData["description"].(string); ok {
		post.Description = desc
	}
	if draft, ok := fmData["draft"].(bool); ok {
		post.Draft = draft
	}
	if dateStr, ok := fmData["date"].(string); ok {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			post.Date = t
		}
	} else if t, ok := fmData["date"].(time.Time); ok {
		post.Date = t
	}

	if tags, ok := fmData["tags"].([]interface{}); ok {
		for _, t := range tags {
			if tag, ok := t.(string); ok {
				post.Tags = append(post.Tags, tag)
			}
		}
	}

	if body == "" {
		body = content
	}
	post.Content = body

	var buf strings.Builder
	if err := p.md.Convert([]byte(body), &buf); err != nil {
		return nil, err
	}
	post.HTML = buf.String()

	return post, nil
}

func (p *Parser) ParseDir(dir string) ([]Post, error) {
	var posts []Post

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		post, err := p.ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("parse post %s: %w", path, err)
		}

		if post.Title == "" {
			post.Title = strings.TrimSuffix(entry.Name(), ".md")
		}

		if post.Date.IsZero() {
			info, err := os.Stat(path)
			if err != nil {
				return nil, fmt.Errorf("stat post %s: %w", path, err)
			}
			post.Date = info.ModTime()
		}

		posts = append(posts, *post)
	}

	sort.Sort(ByDate(posts))

	return posts, nil
}

func Slugify(s string) string {
	s = strings.TrimSuffix(s, ".md")
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
