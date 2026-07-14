package feed

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"time"

	"myblog/config"
	"myblog/content"
)

type RSS struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XmlnsContent string   `xml:"xmlns:content,attr"`
	Channel      Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	LastBuild   time.Time `xml:"lastBuildDate"`
	Items       []Item    `xml:"item"`
}

type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Guid        string   `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
	Description string   `xml:"description"`
	Content     string   `xml:"-"`
	Tags        []string `xml:"category"`
}

func (i Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Item
	if err := e.EncodeElement(struct {
		Alias
		Content string `xml:"content:encoded"`
	}{
		Alias:   Alias(i),
		Content: i.Content,
	}, start); err != nil {
		return err
	}
	return nil
}

type Atom struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Link    []Link   `xml:"link"`
	Updated string   `xml:"updated"`
	Entries []Entry  `xml:"entry"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Entry struct {
	Title      string     `xml:"title"`
	Link       []Link     `xml:"link"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Summary    string     `xml:"summary"`
	Content    string     `xml:"content"`
	Categories []Category `xml:"category"`
}

type Category struct {
	Term string `xml:"term,attr"`
}

func Generate(posts []content.Post, cfg *config.Config, outputDir string) error {
	sort.Sort(content.ByDate(posts))

	siteURL := cfg.AbsURL("")

	rss := RSS{
		Version:      "2.0",
		XmlnsContent: "http://purl.org/rss/1.0/modules/content/",
		Channel: Channel{
			Title:       cfg.Site.Title,
			Link:        siteURL,
			Description: cfg.Site.Description,
			Language:    "en-us",
			LastBuild:   time.Now(),
		},
	}

	var items []Item
	for _, p := range posts {
		if p.Draft {
			continue
		}

		item := Item{
			Title:       p.Title,
			Link:        cfg.AbsURL("posts/" + p.Slug + ".html"),
			Guid:        cfg.AbsURL("posts/" + p.Slug + ".html"),
			PubDate:     p.Date.Format(time.RFC1123),
			Description: p.Description,
			Content:     p.HTML,
			Tags:        p.Tags,
		}
		items = append(items, item)
	}
	rss.Channel.Items = items

	rssData, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	rssPath := filepath.Join(outputDir, "feed.xml")
	if err := os.WriteFile(rssPath, addXMLHeader(rssData), 0644); err != nil {
		return err
	}

	atom := Atom{
		Title:   cfg.Site.Title,
		Link:    []Link{{Rel: "alternate", Href: siteURL}},
		Updated: time.Now().Format(time.RFC3339),
	}

	if len(posts) > 0 {
		atom.Updated = posts[0].Date.Format(time.RFC3339)
	}

	for _, p := range posts {
		if p.Draft {
			continue
		}

		categories := make([]Category, len(p.Tags))
		for i, tag := range p.Tags {
			categories[i] = Category{Term: tag}
		}

		entry := Entry{
			Title:      p.Title,
			Link:       []Link{{Rel: "alternate", Href: cfg.AbsURL("posts/" + p.Slug + ".html")}},
			ID:         cfg.AbsURL("posts/" + p.Slug + ".html"),
			Updated:    p.Date.Format(time.RFC3339),
			Summary:    p.Description,
			Content:    p.HTML,
			Categories: categories,
		}
		atom.Entries = append(atom.Entries, entry)
	}

	atomData, err := xml.MarshalIndent(atom, "", "  ")
	if err != nil {
		return err
	}

	atomPath := filepath.Join(outputDir, "feed.atom")
	return os.WriteFile(atomPath, addXMLHeader(atomData), 0644)
}

func addXMLHeader(data []byte) []byte {
	header := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	return append(header, data...)
}
