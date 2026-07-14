package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myblog/config"
	"myblog/content"
	"myblog/feed"
	"myblog/generator"
	"myblog/search"
)

var (
	configPath string
	buildOnly  bool
	serveFlag  bool
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "Path to config file")
	flag.BoolVar(&buildOnly, "build", true, "Build the site")
	flag.BoolVar(&serveFlag, "serve", false, "Build and serve the site")
	flag.IntVar(&port, "port", 8080, "Port to serve on")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.VisitAll(func(f *flag.Flag) {
			def := ""
			if f.DefValue != "" && f.DefValue != "0" {
				def = fmt.Sprintf(" (default: %s)", f.DefValue)
			}
			fmt.Fprintf(os.Stderr, "  -%s%s\n", f.Name, def)
		})
	}

	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	if buildOnly || serveFlag {
		parser := content.NewParser()
		posts, err := parser.ParseDir(cfg.ResolveProjectPath(cfg.Build.ContentDir))
		if err != nil {
			return fmt.Errorf("failed to parse posts: %w", err)
		}

		gen := generator.New(cfg)
		if err := gen.LoadTemplates(); err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}

		startTime := time.Now()
		log.Printf("Generating blog - %d posts...", len(posts))
		if err := gen.Generate(posts); err != nil {
			return fmt.Errorf("failed to generate blog: %w", err)
		}

		if cfg.Features.Search {
			if err := search.BuildIndex(posts, cfg.OutputPath()); err != nil {
				log.Printf("Warning: failed to build search index: %v", err)
			}
		}

		if cfg.Features.RSS {
			if err := feed.Generate(posts, cfg, cfg.OutputPath()); err != nil {
				log.Printf("Warning: failed to generate feeds: %v", err)
			}
		}

		elapsed := time.Since(startTime)
		if elapsed.Minutes() >= 1 {
			minutes := int(elapsed.Minutes())
			seconds := int(elapsed.Seconds()) % 60
			log.Printf("Blog generated in %d m %d s", minutes, seconds)
		} else if elapsed.Seconds() >= 1 {
			seconds := int(elapsed.Seconds())
			log.Printf("Blog generated in %d s", seconds)
		} else {
			log.Printf("Blog generated in %v", elapsed)
		}
		log.Printf("Blog generated in directory `%s`", cfg.OutputPath())
	}

	if serveFlag {
		return serve(cfg.OutputPath(), port)
	}

	return nil
}

func serve(dir string, port int) error {
	log.Printf("Serving at http://localhost:%d", port)

	fs := http.Dir(dir)
	handler := http.FileServer(fs)

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path != "/" && !strings.HasSuffix(path, "/") {
			fullPath := filepath.Join(dir, path)
			if _, err := os.Stat(fullPath); err != nil {
				if strings.HasSuffix(path, ".html") {
					http.NotFound(w, r)
					return
				}
				if _, err := os.Stat(fullPath + ".html"); err == nil {
					r.URL.Path = path + ".html"
				}
			}
		}

		handler.ServeHTTP(w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), httpHandler)
}
