package vecdb

import "time"

type EmbeddDocument struct {
	// IDMetaKey is the key for the identifing unique mata data
	IDMetaKey string
	// IDMetaValue is the value for the identifing unique mata data
	IDMetaValue string

	URL string

	Modified time.Time
	Title    string
	Document string

	MetaData map[string]any
}
