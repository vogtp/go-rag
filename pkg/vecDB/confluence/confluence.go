package confluence

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/graze/go-throttled"
	"github.com/k3a/html2text"
	"github.com/spf13/viper"
	conflu "github.com/virtomize/confluence-go-api"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"golang.org/x/time/rate"
)

func GetDocuments(ctx context.Context, slog *slog.Logger) (chan vecdb.EmbeddDocument, error) {
	baseURL := viper.GetString(cfg.ConfluenceBaseURL)
	baseURL = strings.TrimRight(baseURL, "/")
	conf := confluence{
		slog:       slog.With("confluence_url", baseURL),
		baseURL:    baseURL,
		out:        make(chan vecdb.EmbeddDocument, 10),
		accessKey:  viper.GetString(cfg.ConfluenceKey),
		rateLimit:  rate.Limit(0.4),
		queryLimit: 100,
		spaces:     viper.GetStringSlice(cfg.ConfluenceSpaces),
	}
	if err := conf.init(); err != nil {
		return nil, err
	}
	go conf.query(ctx)
	return conf.out, nil
}

type confluence struct {
	slog       *slog.Logger
	baseURL    string
	out        chan vecdb.EmbeddDocument
	api        *conflu.API
	accessKey  string
	rateLimit  rate.Limit
	queryLimit int
	spaces     []string
	mu         sync.Mutex
}

func (c *confluence) init() error {
	url := c.getAPIURL()
	api, err := conflu.NewAPI(url, "", c.accessKey)
	if err != nil {
		return err
	}
	if api == nil {
		return fmt.Errorf("confluence api was not created")
	}
	api.Client = throttled.WrapClient(api.Client, rate.NewLimiter(c.rateLimit, 1))
	c.api = api
	c.slog.Info("loaded confluence rest api", "url", url)
	return nil
}

func (c *confluence) getAPIURL() string {
	return c.getURL("/rest/api")
}

func (c *confluence) getURL(ui string) string {
	return fmt.Sprintf("%s%s", c.baseURL, ui)
}

func (c *confluence) query(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer close(c.out)
	for _, s := range c.spaces {
		if ctx.Err() != nil {
			return
		}
		c.querySpace(ctx, s)
	}
}

func (c *confluence) querySpace(ctx context.Context, spaceKey string) {
	slog := c.slog.With("space", spaceKey)
	slog.Info("starting to query space")
	start := 0
	total := 0
	for {
		if ctx.Err() != nil {
			slog.Warn("confluence canceled by context", "err", ctx.Err())
			return
		}
		res, err := c.api.GetContent(conflu.ContentQuery{
			SpaceKey: spaceKey,
			Start:    start,
			Limit:    c.queryLimit,
			Expand:   []string{"space", "body.view", "version", "container", "body.storage", "metadata", "history.lastUpdated"},
		})
		if err != nil {
			log.Fatal(err)
		}
		start += res.Limit
		total += res.Size
		//	fmt.Printf("%s\n ^%s (%v)\n%s (%v)\n", res.Results[0].Title, res.Results[0].Links.WebUI, len(res.Results[0].Body.Storage.Value), res.Results[res.Size-1].Title, len(res.Results[res.Size-1].Body.Storage.Value))

		for _, d := range res.Results {
			if ctx.Err() != nil {
				slog.Warn("confluence canceled by context", "err", ctx.Err())
				return
			}
			slog.Debug("processing confluence document", "title", d.Title)
			txt := html2text.HTML2Text(d.Body.View.Value)
			doc := vecdb.EmbeddDocument{
				Title:       d.Title,
				URL:         c.getURL(d.Links.WebUI),
				Document:    txt,
				IDMetaKey:   vecdb.MetaURL,
				IDMetaValue: d.Links.WebUI,
			}
			// 2016-05-30T16:14:07.787+02:00
			if t, err := time.Parse(time.RFC3339Nano, d.History.LastUpdated.When); err == nil {
				doc.Modified = t
			} else {
				c.slog.Error("Cannot parse time of confluence page", "time", t.String(), "err", err, "title", d.Title, "url", d.Links.WebUI)
				continue
			}
			c.out <- doc
		}

		slog.Info("confluence query batch done", "start", res.Start, "size", res.Size, "result_len", len(res.Results), "limit", res.Limit, "total", total)
		if res.Limit != res.Size {
			break
		}
	}
}

// func (c *confluence) getAllSpaces() {
// 	panic("Not implemented")
// 	// spaces, err := api.GetAllSpaces(conflu.AllSpacesQuery{
// 	// 	Type:   "global",
// 	// 	Start:  0,
// 	// 	Limit:  999,
// 	// 	Expand: []string{"space", "body.view", "version", "container"},
// 	// })
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// for _, space := range spaces.Results {
// }
