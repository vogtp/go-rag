package scraper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

const (
	DefualtMaxDept   = 1
	DefualtParallels = 2
	DefualtDelay     = 3
	DefualtAsync     = true
)

var ErrScrapingFailed = errors.New("scraper could not read URL, or scraping is not allowed for provided URL")

type Scraper struct {
	MaxDepth  int
	Parallels int
	Delay     int64
	Blacklist []string
	Async     bool
}

//var _ documentloaders.Loader = Scraper{}

// New creates a new instance of Scraper with the provided options.
//
// The options parameter is a variadic argument allowing the user to specify
// custom configuration options for the Scraper. These options can be
// functions that modify the Scraper's properties.
//
// The function returns a pointer to a Scraper instance and an error. The
// error value is nil if the Scraper is created successfully.
func New(options ...Options) (*Scraper, error) {
	scraper := &Scraper{
		MaxDepth:  DefualtMaxDept,
		Parallels: DefualtParallels,
		Delay:     int64(DefualtDelay),
		Async:     DefualtAsync,
		Blacklist: []string{
			"login",
			"signup",
			"signin",
			"register",
			"logout",
			"download",
			"redirect",
		},
	}

	for _, opt := range options {
		opt(scraper)
	}

	return scraper, nil
}

// Name returns the name of the Scraper.
//
// No parameters.
// Returns a string.
func (s Scraper) Name() string {
	return "Web Scraper"
}

// Description returns the description of the Go function.
//
// There are no parameters.
// It returns a string.
func (s Scraper) Description() string {
	return `
		Web Scraper will scan a url and return the content of the web page.
		Input should be a working url.
	`
}

// Call scrapes a website and returns the site data.
//
// The function takes a context.Context object for managing the execution
// context and a string input representing the URL of the website to be scraped.
// It returns a string containing the scraped data and an error if any.
//
//nolint:all
func (s Scraper) Call(ctx context.Context, input string, docsOutput chan vecdb.EmbeddDocument) error {
	defer close(docsOutput)
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	c := colly.NewCollector(
		colly.MaxDepth(s.MaxDepth),
		colly.Async(s.Async),
		// colly.MaxDepth(1),
		//colly.Debugger(&debug.LogDebugger{Prefix: "XXX"}),
	)

	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: s.Parallels,
		Delay:       time.Duration(s.Delay) * time.Second,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	var siteData stringBuilder
	homePageLinks := make(map[string]bool)
	scrapedLinks := make(map[string]bool)
	scrapedLinksMutex := sync.RWMutex{}

	c.OnRequest(func(r *colly.Request) {
		if ctx.Err() != nil {
			r.Abort()
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		if ctx.Err() != nil {
			slog.Info("Context canceled", "err", ctx.Err())
			return
		}
		currentURL := e.Request.URL.String()
		siteData.Reset()
		// Only process the page if it hasn't been visited yet
		scrapedLinksMutex.Lock()
		if !scrapedLinks[currentURL] {
			scrapedLinks[currentURL] = true
			scrapedLinksMutex.Unlock()

			siteData.WriteString("\n\nPage URL: " + currentURL)
			slog.Info("Processing Page", "URL", currentURL)
			title := e.ChildText("title")
			if title != "" {
				siteData.WriteString("\nPage Title: " + title)
			}

			description := e.ChildAttr("meta[name=description]", "content")
			if description != "" {
				siteData.WriteString("\nPage Description: " + description)
			}
			doc := vecdb.EmbeddDocument{
				IDMetaKey:   vecdb.MetaURL,
				IDMetaValue: currentURL,
				Title:       title,
				Modified: time.Now(),
				MetaData:    make(map[string]any),
			}
			if len(description) > 0 {
				doc.MetaData["description"] = description
			}

			siteData.WriteString("\nHeaders:")
			e.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
				siteData.WriteString("\n" + el.Text)
			})

			siteData.WriteString("\nContent:")
			e.ForEach("p", func(_ int, el *colly.HTMLElement) {
				siteData.WriteString("\n" + el.Text)
			})

			doc.Document = siteData.String()
			slog.Debug("Sending document to embedder")
			docsOutput <- doc

			if currentURL == input {
				e.ForEach("a", func(_ int, el *colly.HTMLElement) {
					link := el.Attr("href")
					if link != "" && !homePageLinks[link] {
						homePageLinks[link] = true
						slog.InfoContext(ctx, "Found link", "link", link)
					}
				})
			}

		} else {
			scrapedLinksMutex.Unlock()
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if ctx.Err() != nil {
			slog.Info("Context canceled", "err", ctx.Err())
			return
		}
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)

		// Parse the link to get the hostname
		u, err := url.Parse(absoluteLink)
		if err != nil {
			// Handle the error appropriately
			return
		}

		// Check if the link's hostname matches the current request's hostname
		if u.Hostname() != e.Request.URL.Hostname() {
			return
		}

		// Check for redundant pages
		for _, item := range s.Blacklist {
			if strings.Contains(u.Path, item) {
				return
			}
		}

		// Normalize the path to treat '/' and '/index.html' as the same path
		if u.Path == "/index.html" || u.Path == "" {
			u.Path = "/"
		}

		// Only visit the page if it hasn't been visited yet
		scrapedLinksMutex.RLock()
		if !scrapedLinks[u.String()] {
			scrapedLinksMutex.RUnlock()
			err := c.Visit(u.String())
			if err != nil {
				if errors.Is(err, colly.ErrAlreadyVisited) {
					slog.DebugContext(ctx, "Not following link", "link", link, "err", err)
				} else {
					slog.WarnContext(ctx, "Error following link", "link", link, "err", err)
				}
			}
		} else {
			scrapedLinksMutex.RUnlock()
		}
	})

	err = c.Visit(input)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrScrapingFailed, err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.Wait()
	}
	slog.InfoContext(ctx, "Finished scrapping", "url", input)
	return nil
}
