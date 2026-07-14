package search

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"myblog/content"
)

type SearchIndex struct {
	Posts []SearchPost `json:"posts"`
}

type SearchPost struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
}

func BuildIndex(posts []content.Post, outputDir string) error {
	var idx []SearchPost

	for _, p := range posts {
		if p.Draft {
			continue
		}

		idx = append(idx, SearchPost{
			Slug:        p.Slug,
			Title:       p.Title,
			Date:        p.Date.Format("2006-01-02"),
			Tags:        p.Tags,
			Description: p.Description,
			Content:     stripHTML(p.HTML),
		})
	}

	data, err := json.MarshalIndent(SearchIndex{Posts: idx}, "", "  ")
	if err != nil {
		return err
	}

	outPath := filepath.Join(outputDir, "search.json")
	return os.WriteFile(outPath, data, 0644)
}

func stripHTML(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	inTag := false
	inSpace := true // Start with true to avoid leading space

	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if inTag {
			continue
		}

		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !inSpace {
				buf.WriteByte(' ')
				inSpace = true
			}
		} else {
			buf.WriteRune(r)
			inSpace = false
		}
	}

	res := buf.String()
	if len(res) > 0 && res[len(res)-1] == ' ' {
		res = res[:len(res)-1]
	}
	return res
}

type Filter struct {
	Posts       []content.Post
	searchTexts []string
}

func NewFilter(posts []content.Post) *Filter {
	var published []content.Post
	for _, p := range posts {
		if !p.Draft {
			published = append(published, p)
		}
	}
	sort.Sort(content.ByDate(published))

	searchTexts := make([]string, len(published))
	for i, p := range published {
		searchTexts[i] = strings.ToLower(p.Title + " " + p.Description + " " + strings.Join(p.Tags, " "))
	}

	return &Filter{
		Posts:       published,
		searchTexts: searchTexts,
	}
}

func (f *Filter) Search(query string) []content.Post {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return f.Posts
	}

	var results []content.Post
	terms := strings.Fields(query)

	for i, p := range f.Posts {
		text := f.searchTexts[i]
		match := true
		for _, term := range terms {
			if !strings.Contains(text, term) {
				match = false
				break
			}
		}
		if match {
			results = append(results, p)
		}
	}

	return results
}
