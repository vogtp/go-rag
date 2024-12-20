package scraper

import (
	"strings"
)

type stringBuilder struct {
	strings.Builder
}

func (b *stringBuilder) WriteString(s string) (int, error) {
	//fmt.Print(s)
	return b.Builder.WriteString(s)
}
