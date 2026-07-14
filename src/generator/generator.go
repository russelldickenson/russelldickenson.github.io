package generator

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"myblog/config"
	"myblog/content"
)

type Generator struct {
	cfg       *config.Config
	posts     []content.Post
	templates *template.Template
	funcMap   template.FuncMap
}

type PageData struct {
	Site          *config.Config
	Posts         []content.Post
	Post          *content.Post
	Tags          map[string]int
	RecentPosts   []content.Post
	Content       template.HTML
	SidebarLeft   bool
	SearchEnabled bool
	TagName       string
	CurrentPage   int
	TotalPages    int
	PrevLink      string
	NextLink      string
	HasPrev       bool
	HasNext       bool
	CanonicalPath string
}

func New(cfg *config.Config) *Generator {
	g := &Generator{
		cfg:     cfg,
		funcMap: template.FuncMap{},
	}

	g.funcMap["formatDate"] = func(d time.Time) string {
		return d.Format("January 2, 2006")
	}
	g.funcMap["formatISO"] = func(d time.Time) string {
		return d.Format(time.RFC3339)
	}
	g.funcMap["json"] = func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	g.funcMap["html"] = func(s string) template.HTML {
		return template.HTML(s)
	}
	g.funcMap["escape"] = html.EscapeString
	g.funcMap["urlize"] = func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
	}
	g.funcMap["now"] = func() time.Time {
		return time.Now()
	}
	g.funcMap["hasPost"] = func(p *content.Post) bool {
		if p == nil {
			return false
		}
		return p.Title != ""
	}
	g.funcMap["baseURL"] = func() string {
		return g.cfg.BaseURL()
	}
	g.funcMap["relURL"] = func(resourcePath string) string {
		return g.cfg.RelURL(resourcePath)
	}
	g.funcMap["absURL"] = func(resourcePath string) string {
		return g.cfg.AbsURL(resourcePath)
	}

	return g
}

func (g *Generator) LoadTemplates() error {
	themePath := g.cfg.ResolveProjectPath("themes")

	pattern := themePath + "/*.html"
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no templates found in %s", themePath)
	}

	tmpl := template.Must(template.New("layout").Funcs(g.funcMap).Parse(""))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read template %s: %w", file, err)
		}
		tmpl, err = tmpl.Parse(string(data))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", file, err)
		}
	}

	g.templates = tmpl
	return nil
}

func (g *Generator) Generate(posts []content.Post) error {
	g.posts = posts

	outputDir := g.cfg.OutputPath()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if err := g.copyStatic(); err != nil {
		return err
	}

	data := g.buildPageData()

	emptyPost := content.Post{}
	data.Post = &emptyPost

	// Pagination logic
	published := data.Posts
	perPage := g.cfg.Build.PostsPerPage
	totalPages := (len(published) + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}

	for i := 0; i < totalPages; i++ {
		start := i * perPage
		end := start + perPage
		if end > len(published) {
			end = len(published)
		}

		pagePosts := published[start:end]
		pageData := data
		pageData.Posts = pagePosts
		pageData.CurrentPage = i + 1
		pageData.TotalPages = totalPages

		if err := g.generateIndex(pageData, i+1, totalPages); err != nil {
			return err
		}
	}

	for _, post := range posts {
		if post.Draft {
			continue
		}
		if err := g.generatePost(post, data); err != nil {
			return fmt.Errorf("generate post %s: %w", post.Slug, err)
		}
	}

	if err := g.generateTagPages(data); err != nil {
		return err
	}

	return nil
}

func (g *Generator) copyStatic() error {
	staticSrc := g.cfg.ResolveProjectPath("static")
	staticDst := filepath.Join(g.cfg.OutputPath(), "static")

	if err := os.MkdirAll(staticDst, 0755); err != nil {
		return err
	}

	// Copy static files
	if _, err := os.Stat(staticSrc); err == nil {
		if err := copyDirExcluding(staticSrc, staticDst, []string{"style.css"}); err != nil {
			return err
		}
	}

	// Copy unified style CSS
	styleCSSPath := filepath.Join(g.cfg.ResolveProjectPath("themes"), "style.css")
	themesDst := filepath.Join(g.cfg.OutputPath(), "themes")
	if err := os.MkdirAll(themesDst, 0755); err != nil {
		return err
	}
	if data, err := os.ReadFile(styleCSSPath); err == nil {
		styleDst := filepath.Join(themesDst, "style.css")
		if err := os.WriteFile(styleDst, data, 0644); err != nil {
			return fmt.Errorf("failed to write style CSS: %w", err)
		}
	} else {
		return fmt.Errorf("style CSS not found at %s: %w", styleCSSPath, err)
	}

	return nil
}

func (g *Generator) buildPageData() PageData {
	var published []content.Post
	for _, p := range g.posts {
		if !p.Draft {
			published = append(published, p)
		}
	}

	sorted := make([]content.Post, len(published))
	copy(sorted, published)
	sort.Sort(content.ByDate(sorted))

	var recent []content.Post
	if len(sorted) > 0 {
		count := g.cfg.Sidebar.RecentCount
		if count > len(sorted) {
			count = len(sorted)
		}
		recent = sorted[:count]
	}

	tags := make(map[string]int)
	for _, p := range published {
		for _, tag := range p.Tags {
			tag = strings.ToLower(tag)
			tags[tag]++
		}
	}

	sidebarLeft := g.cfg.Sidebar.Position == "left"

	return PageData{
		Site:          g.cfg,
		Posts:         published,
		Tags:          tags,
		RecentPosts:   recent,
		SidebarLeft:   sidebarLeft,
		SearchEnabled: g.cfg.Features.Search,
		CurrentPage:   1,
		TotalPages:    1,
		PrevLink:      "",
		NextLink:      "",
		HasPrev:       false,
		HasNext:       false,
	}
}

func (g *Generator) generateIndex(data PageData, page int, totalPages int) error {
	data.HasPrev = page > 1
	data.HasNext = page < totalPages

	if data.HasPrev {
		if page == 2 {
			data.PrevLink = g.cfg.RelURL("")
		} else {
			data.PrevLink = g.cfg.RelURL(fmt.Sprintf("page/%d/index.html", page-1))
		}
	}

	if data.HasNext {
		data.NextLink = g.cfg.RelURL(fmt.Sprintf("page/%d/index.html", page+1))
	}

	if page == 1 {
		data.CanonicalPath = ""
	} else {
		data.CanonicalPath = fmt.Sprintf("page/%d/index.html", page)
	}

	output := "index.html"
	if page > 1 {
		output = fmt.Sprintf("page/%d/index.html", page)
	}

	return g.renderTemplate("layout", data, output)
}

func (g *Generator) generatePost(post content.Post, data PageData) error {
	data.Post = &post
	data.Content = template.HTML(post.HTML)
	data.CanonicalPath = filepath.Join("posts", post.Slug+".html")

	filename := filepath.Join("posts", post.Slug+".html")

	return g.renderTemplate("post", data, filename)
}

func (g *Generator) generateTagPages(data PageData) error {
	tagsDir := filepath.Join(g.cfg.OutputPath(), "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	for tag := range data.Tags {
		var tagPosts []content.Post
		for _, post := range data.Posts {
			if containsTag(post.Tags, tag) {
				tagPosts = append(tagPosts, post)
			}
		}

		tagData := data
		tagData.Posts = tagPosts
		tagData.TagName = tag
		tagSlug := content.Slugify(tag)
		tagData.CanonicalPath = filepath.Join("tags", tagSlug+".html")

		tagFile := tagSlug + ".html"
		if err := g.renderTemplate("tag", tagData, "tags/"+tagFile); err != nil {
			return fmt.Errorf("generate tag page %s: %w", tag, err)
		}
	}

	return nil
}

func containsTag(tags []string, tag string) bool {
	tagLower := strings.ToLower(tag)
	for _, t := range tags {
		if strings.ToLower(t) == tagLower {
			return true
		}
	}
	return false
}

func (g *Generator) renderTemplate(name string, data interface{}, output string) error {
	tmpl := g.templates.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("execute %s: %w", name, err)
	}

	outPath := filepath.Join(g.cfg.OutputPath(), output)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

func copyDirExcluding(src, dst string, exclude []string) error {
	excludeMap := make(map[string]bool)
	for _, f := range exclude {
		excludeMap[f] = true
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		if excludeMap[rel] {
			return nil
		}

		dstPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
