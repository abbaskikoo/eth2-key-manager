package handler

import "github.com/bloxapp/KeyVault/cli/util/printer"

// Seed contains handler functions of the CLI commands related to portfolio seed.
type Seed struct {
	printer printer.Printer
}

// New is the constructor of Seed handler.
func New(printer printer.Printer) *Seed {
	return &Seed{
		printer: printer,
	}
}