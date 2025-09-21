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

	var allSymbols []*domain.Symbol

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

		allSymbols = append(allSymbols, symbols...)
	}

	kind := map[domain.SymbolKind]string{
		domain.FunctionDefinition:  "definition",
		domain.FunctionDeclaration: "declaration",
		domain.FunctionCall:        "call",
	}

	groups := domain.NewGroups(allSymbols)
	for _, group := range groups {
		if len(group.Definitions) == 0 {
			continue
		}

		s.printer.Print("%s\n", group.Definitions[0].Name)

		for _, definition := range group.Definitions {
			s.printer.Print(" - definition %s %s:%d\n", kind[definition.Kind], definition.Path, definition.Line)
		}

		for _, declaration := range group.Declarations {
			s.printer.Print(" - declaration %s:%d\n", declaration.Path, declaration.Line)
		}

		for _, call := range group.Calls {
			s.printer.Print(" - call %s:%d\n", call.Path, call.Line)
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
