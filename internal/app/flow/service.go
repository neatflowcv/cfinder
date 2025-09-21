package flow

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/neatflowcv/cfinder/internal/pkg/filesystem"
	"github.com/neatflowcv/cfinder/internal/pkg/finder"
	"github.com/neatflowcv/cfinder/internal/pkg/parser"
	"github.com/neatflowcv/cfinder/internal/pkg/printer"
)

type Service struct {
	filesystem filesystem.Filesystem
	printer    printer.Printer
}

func NewService(
	filesystem filesystem.Filesystem,
	printer printer.Printer,
) *Service {
	return &Service{
		filesystem: filesystem,
		printer:    printer,
	}
}

func (s *Service) FindSymbol(ctx context.Context, dir string, symbol string) error {
	files, err := s.filesystem.ListFiles(ctx, dir)
	if err != nil {
		return fmt.Errorf("ListFiles: %w", err)
	}

	var filteredFiles []string

	for _, file := range files {
		if !s.isCodeFile(file) {
			continue
		}

		filteredFiles = append(filteredFiles, file)
	}

	for _, path := range filteredFiles {
		file, err := s.filesystem.OpenFile(ctx, path)
		if err != nil {
			return fmt.Errorf("OpenFile: %w", err)
		}

		finder := finder.NewFinder(file)
		symbols := finder.FindSymbol(symbol)

		_ = file.Close()

		if len(symbols) > 0 {
			for _, symbol := range symbols {
				s.printer.Print("%s:%d\n", path, symbol.Line)
			}
		}
	}

	return nil
}

func (s *Service) ListSymbols(ctx context.Context, dir string) error {
	files, err := s.filesystem.ListFiles(ctx, dir)
	if err != nil {
		return fmt.Errorf("ListFiles: %w", err)
	}

	var filteredFiles []string

	for _, file := range files {
		if !s.isCodeFile(file) {
			continue
		}

		filteredFiles = append(filteredFiles, file)
	}

	for _, path := range filteredFiles {
		file, err := s.filesystem.OpenFile(ctx, path)
		if err != nil {
			return fmt.Errorf("OpenFile: %w", err)
		}

		parser := parser.NewParser()

		symbols, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("ParseFile: %w", err)
		}

		_ = file.Close()

		for _, symbol := range symbols {
			s.printer.Print("%s:%d %s\n", path, symbol.Line, symbol.Name)
		}
	}

	return nil
}

func (s *Service) isCodeFile(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))

	return ext == ".c" ||
		ext == ".cpp" ||
		ext == ".cc" ||
		ext == ".h" ||
		ext == ".hpp"
}
