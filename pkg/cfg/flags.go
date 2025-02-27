package cfg

import (
	"fmt"
	"time"

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

	// OllamaHosts is an URL
	ollamaHosts = "ollama.hosts"

	// ChromaUrl the URL where chroma can be reached
	ChromaUrl = "chroma.url"
	// ChromaPort is the port chroma should be started on (0: disable)
	ChromaPort = "chroma.port"
	// ChromaContainer chroma container to pull
	ChromaContainer = "chroma.container"

	// ConfluenceKey is the confluence access token
	ConfluenceKey = "confluence.key"
	// ConfluenceBaseURL is the base URL of the confluence instance
	ConfluenceBaseURL = "confluence.baseURL"
	// ConfluenceSpaces defines the spaces to scrap
	ConfluenceSpaces = "confluence.spaces"
	//VecDBUpdateIntervall is the intervall the vectorDB is updated
	VecDBUpdateIntervall = "vecdb.update_intervall"
)

func init() {
	pflag.Bool(CfgSave, false, "Should the configs be written to file periodically")
	pflag.String(CfgFile, fmt.Sprintf("%s.yml", appName), "File with the config to load")
	pflag.String(LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	pflag.Bool(LogSource, false, "Log the source line")
	pflag.Bool(LogJson, false, "Log in json")

	pflag.String(WebListen, ":8080", "Address the webserver should listen on")

	pflag.String(ModelDefault, "llama3.2-vision", "The default model when no model is given")
	pflag.String(ModelEmbedding, "mxbai-embed-large", "The default model used for embeddings")

	pflag.String(ConfluenceKey, "", "The confluence access token")
	pflag.String(ConfluenceBaseURL, "", "The confluence access token")
	pflag.StringSlice(ConfluenceSpaces, nil, "The confluence spaces to scrap")
	pflag.Duration(VecDBUpdateIntervall, 24*time.Hour, "the intervall the vectorDB is updated")
	pflag.String(ChromaUrl, "http://localhost:8000", "the URL where chroma can be reached")
	pflag.Int(ChromaPort, 8000, "the port chroma should be started on (0: disable)")
	pflag.String(ChromaContainer, "chromadb/chroma:0.5.23", "chroma container to pull")
}
