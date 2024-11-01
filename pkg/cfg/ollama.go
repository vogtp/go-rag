package cfg

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	ollamaAPI "github.com/ollama/ollama/api"
	"github.com/spf13/viper"
)

var _ollamaHost string

func GetOllamaHost(ctx context.Context) string {
	if _ollamaHost != "" {
		return _ollamaHost
	}
	for _, o := range viper.GetStringSlice(OllamaHosts) {
		u, err := url.Parse(o)
		if err != nil {
			slog.Warn("Cannot parse ollama url", "url", o, "err", err)
			continue
		}
		c := ollamaAPI.NewClient(u, http.DefaultClient)
		if err := c.Heartbeat(ctx); err != nil {
			slog.Warn("Cannot connect to ollama", "url", o, "err", err)
		}
		_ollamaHost = o
		break
	}
	return _ollamaHost
}
