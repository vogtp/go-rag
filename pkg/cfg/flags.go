package cfg

import (
	"fmt"

	"github.com/spf13/pflag"
)

const (
	// CfgFile
	CfgFile = "config.file"
	// CfgSave triggers periodic config saves
	CfgSave = "config.save"

	// WebListen set the address the webserver should listen on
	WebListen = "web.listen"

	// LogLevel error warn info debug
	LogLevel = "log.level"
	// LogSource should we log the source
	LogSource = "log.source"
	// LogJson log in json
	LogJson = "log.json"

	// CacheDir the subdirectory to use as cache, subdir of modeldir if it does not start with /
	CacheDir = "cache.dir"

	// ModelDefault is the default model when no model is given
	ModelDefault = "model.default"
	// ModelEmbedding is the default model used for embeddings
	ModelEmbedding = "model.embedding"

	// OllamaHosts is host:port
	OllamaHosts = "ollama.hosts"

	// ChromaUrl the URL where chroma can be reached
	ChromaUrl = "chroma.url"

	// ConfluenceKey is the confluence access token
	ConfluenceKey = "confluence.key"
	// ConfluenceBaseURL is the base URL of the confluence instance
	ConfluenceBaseURL = "confluence.baseURL"
	// ConfluenceSpaces defines the spaces to scrap
	ConfluenceSpaces = "confluence.spaces"
)

func init() {
	pflag.Bool(CfgSave, false, "Should the configs be written to file periodically")
	pflag.String(CfgFile, fmt.Sprintf("%s.yml", APP_NAME), "File with the config to load")
	pflag.String(LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	pflag.Bool(LogSource, false, "Log the source line")
	pflag.Bool(LogJson, false, "Log in json")

	pflag.String(ModelDefault, "llama3.1", "The default model when no model is given")
	pflag.String(ModelEmbedding, "mxbai-embed-large", "The default model used for embeddings")

	pflag.String(ConfluenceKey, "", "The confluence access token")
	pflag.String(ConfluenceBaseURL, "", "The confluence access token")
	pflag.StringSlice(ConfluenceSpaces, nil, "The confluence spaces to scrap")

	pflag.String(ChromaUrl, "http://localhost:8000", "the URL where chroma can be reached")
}
