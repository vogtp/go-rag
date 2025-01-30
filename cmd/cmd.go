package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
)

func New() *cobra.Command {

	rootFlags()
	chromaFlags()
	testFlags()

	cfg.Parse()
	logger.New()

	addRoot()

	addVecDB()
	addOllama()
	addchroma()
	addWeb()
	addTest()

	return rootCmd
}
