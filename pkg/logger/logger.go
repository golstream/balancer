package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Init() {
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return slog.Attr{}
			}
			source.File = filepath.Base(source.File)
		}
		return a
	}

	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replacer})),
	)
}
