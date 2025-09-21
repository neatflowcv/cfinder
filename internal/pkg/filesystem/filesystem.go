package filesystem

import (
	"context"
	"io"
)

type Filesystem interface {
	ListFiles(ctx context.Context, dir string, excludes []string) ([]string, error)
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
}
