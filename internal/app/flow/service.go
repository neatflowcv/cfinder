package flow

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/neatflowcv/cfinder/internal/pkg/domain"
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
	files, err := s.filesystem.ListFiles(ctx, dir, nil)
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
		symbols := finder.FindSymbol(path, symbol)

		_ = file.Close()

		if len(symbols) > 0 {
			for _, symbol := range symbols {
				s.printer.Print("%s:%d\n", symbol.Path, symbol.Line)
			}
		}
	}

	return nil
}

func (s *Service) ListSymbols(ctx context.Context, dir string, excludes []string) error { //nolint:cyclop,funlen
	files, err := s.filesystem.ListFiles(ctx, dir, excludes)
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

		symbols, err := parser.ParseFile(path, file)
		if err != nil {
			return fmt.Errorf("ParseFile: %w", err)
		}

		_ = file.Close()

		kind := map[domain.SymbolKind]string{
			domain.FunctionDefinition:  "definition",
			domain.FunctionDeclaration: "declaration",
			domain.FunctionCall:        "call",
		}

		groups := domain.NewGroups(symbols)

		for _, group := range groups {
			if group.Definition == nil {
				continue
			}

			s.printer.Print("%s\n", group.Definition.Name)

			if len(group.Calls) == 0 {
				s.printer.Print(" - no call %s %s:%d\n", kind[group.Definition.Kind], group.Definition.Path, group.Definition.Line)

				continue
			}

			s.printer.Print(" - %s %s:%d\n", kind[group.Definition.Kind], group.Definition.Path, group.Definition.Line)

			if group.Declaration != nil {
				s.printer.Print(" - declaration %s:%d\n", group.Declaration.Path, group.Declaration.Line)
			}

			for _, call := range group.Calls {
				s.printer.Print(" - call %s:%d\n", call.Path, call.Line)
			}
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
