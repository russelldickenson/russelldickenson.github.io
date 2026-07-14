package search

import (
	"reflect"
	"testing"
	"time"

	"myblog/content"
)

func TestStripHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"   ", ""},
		{"<p>Hello World</p>", "Hello World"},
		{"  <p>Hello  \n  World</p>  ", "Hello World"},
		{"<div class=\"test\">Nested <span>Text</span> here.</div>", "Nested Text here."},
		{"No HTML tags at all.", "No HTML tags at all."},
		{"<invalid>tag", "tag"},
		{"word<tag>word", "wordword"},
		{"word <tag> word", "word word"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := stripHTML(tc.input)
			if got != tc.expected {
				t.Errorf("stripHTML(%q) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestFilterSearch(t *testing.T) {
	posts := []content.Post{
		{
			Title:       "Introduction to Go",
			Description: "An introductory post about programming in Go.",
			Tags:        []string{"Go", "Programming"},
			Date:        time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
			Draft:       false,
		},
		{
			Title:       "Optimizing Web Applications",
			Description: "Tips and tricks for making websites load faster.",
			Tags:        []string{"Web", "Performance", "HTML"},
			Date:        time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC),
			Draft:       false,
		},
		{
			Title:       "Draft Post Title",
			Description: "This should not be searchable.",
			Tags:        []string{"Draft"},
			Date:        time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC),
			Draft:       true,
		},
	}

	filter := NewFilter(posts)

	// Test that draft is excluded and posts are sorted by date descending (Optimizing Web Applications first)
	if len(filter.Posts) != 2 {
		t.Errorf("expected 2 published posts, got %d", len(filter.Posts))
	}
	if filter.Posts[0].Title != "Optimizing Web Applications" {
		t.Errorf("expected latest post first, got %q", filter.Posts[0].Title)
	}

	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "empty query returns all",
			query:    "",
			expected: []string{"Optimizing Web Applications", "Introduction to Go"},
		},
		{
			name:     "query matching title",
			query:    "introduction",
			expected: []string{"Introduction to Go"},
		},
		{
			name:     "query matching description",
			query:    "faster",
			expected: []string{"Optimizing Web Applications"},
		},
		{
			name:     "query matching tags",
			query:    "programming",
			expected: []string{"Introduction to Go"},
		},
		{
			name:     "case insensitivity",
			query:    "gO",
			expected: []string{"Introduction to Go"},
		},
		{
			name:     "multiple terms",
			query:    "web performance",
			expected: []string{"Optimizing Web Applications"},
		},
		{
			name:     "no match",
			query:    "rust",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := filter.Search(tc.query)
			var got []string
			for _, r := range results {
				got = append(got, r.Title)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Search(%q) = %v; want %v", tc.query, got, tc.expected)
			}
		})
	}
}
