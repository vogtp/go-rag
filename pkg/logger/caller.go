package logger

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

var (
	srcPrefix    string
	srcPrefixLen int
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return
	}
	idx := strings.Index(filename, `rag/pkg/logger/caller.go`)
	if idx > len(filename) || idx < 0 {
		return
	}
	srcPrefix = filename[:idx]
	srcPrefixLen = len(srcPrefix)
}

func trimPackagePath(s string) string {
	if len(s) > len(srcPrefix) && strings.HasPrefix(s, srcPrefix) {
		return s[srcPrefixLen:]
	}
	return s
}

type source struct {
	*slog.Source
	Src string `json:"source"`
}

func ProcessSourceField(attr slog.Attr, json bool) slog.Attr {
	src, ok := attr.Value.Any().(*slog.Source)
	if !ok {
		return attr
	}
	src.File = trimPackagePath(src.File)
	line := fmt.Sprintf("%s:%d", src.File, src.Line)
	if json {
		ret := &source{
			Source: src,
			Src:    line,
		}
		attr.Value = slog.AnyValue(ret)
	} else {
		attr.Value = slog.StringValue(line)
	}
	return attr
}

func Caller(skip int) string {
	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", trimPackagePath(filename), line)
}
