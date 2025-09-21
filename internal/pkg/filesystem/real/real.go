package realfilesystem

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/neatflowcv/cfinder/internal/pkg/filesystem"
)

var _ filesystem.Filesystem = (*Filesystem)(nil)

type Filesystem struct{}

func NewFilesystem() *Filesystem {
	return &Filesystem{}
}

func (f *Filesystem) ListFiles(ctx context.Context, dir string, excludes []string) ([]string, error) {
	expandedDir, err := expandHomeDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string

	err = filepath.WalkDir(expandedDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() && slices.Contains(excludes, path) {
			return filepath.SkipDir
		}

		if slices.Contains(excludes, path) {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ListFiles: %w", err)
	}

	return files, nil
}

func (f *Filesystem) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	newVar, err := os.Open(path) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("OpenFile: %w", err)
	}

	return newVar, nil
}

func expandHomeDir(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("expandHomeDir: %w", err)
	}

	if path == "~" {
		return homeDir, nil
	}

	return filepath.Join(homeDir, path[1:]), nil
}
