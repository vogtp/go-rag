package vecdb

import "time"

type EmbeddDocument struct {
	IDMetaKey   string
	IDMetaValue string
	URL         string

	Modified time.Time
	Title    string
	Document string

	MetaData map[string]any
}
