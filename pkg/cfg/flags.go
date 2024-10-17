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
)

func init() {
	pflag.Bool(CfgSave, false, "Should the configs be written to file periodically")
	pflag.String(CfgFile, fmt.Sprintf("%s.yml", APP_NAME), "File with the config to load")
	pflag.String(LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	pflag.Bool(LogSource, false, "Log the source line")
	pflag.Bool(LogJson, false, "Log in json")

	pflag.String(ModelDefault, "llama3.1", "ModelDefault is the default model when no model is given")
	pflag.String(ModelEmbedding, "llama3.1", "ModelEmbedding is the default model used for embeddings")
}
