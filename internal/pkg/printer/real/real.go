package realprinter

import (
	"fmt"
	"log"

	"github.com/neatflowcv/cfinder/internal/pkg/printer"
)

var _ printer.Printer = (*Printer)(nil)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(format string, args ...any) {
	_, err := fmt.Printf(format, args...) //nolint:forbidigo
	if err != nil {
		log.Fatal(err)
	}
}
