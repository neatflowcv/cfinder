package main

import (
	"context"
	"log"
	"os"

	"github.com/neatflowcv/cfinder/internal/app/flow"
	realfilesystem "github.com/neatflowcv/cfinder/internal/pkg/filesystem/real"
	realprinter "github.com/neatflowcv/cfinder/internal/pkg/printer/real"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
)

// isCppFile checks if the file has a C/C++ extension
// func isCppFile(filename string) bool {
// 	ext := strings.ToLower(filepath.Ext(filename))
// 	return ext == ".c" || ext == ".cpp" || ext == ".cc" || ext == ".h" || ext == ".hpp"
// }

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("version", version)

	service := flow.NewService(
		realfilesystem.NewFilesystem(),
		realprinter.NewPrinter(),
	)

	cmd := cli.Command{ //nolint:exhaustruct
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "dir",
				Value: ".",
				Usage: "directory path",
			},
			&cli.StringSliceFlag{ //nolint:exhaustruct
				Name:  "excludes",
				Value: nil,
				Usage: "exclude",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "find",
				Action: func(ctx context.Context, c *cli.Command) error {
					dir := c.String("dir")
					symbol := c.Args().First()
					if symbol == "" {
						log.Fatal("symbol is required")
					}

					return service.FindSymbol(ctx, dir, symbol)
				},
			},
			{
				Name: "symbols",
				Action: func(ctx context.Context, c *cli.Command) error {
					dir := c.String("dir")
					excludes := c.StringSlice("excludes")

					return service.ListSymbols(ctx, dir, excludes)
				},
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
