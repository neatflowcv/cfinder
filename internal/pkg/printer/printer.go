package printer

type Printer interface {
	Print(format string, args ...any)
}
