package main

import (
	"errors"
	"os"

	"github.com/jessevdk/go-flags"
)

func handleFlagsError(err error) {
	if err != nil {
		var e *flags.Error
		if errors.As(err, &e) {
			switch {
			case errors.Is(e, flags.ErrHelp):
				os.Exit(0)
			case errors.Is(e, flags.ErrRequired):
				os.Exit(1)
			}
		} else {
			// probably a command error
			// go-flags already prints it for us
			os.Exit(1)
		}
	}
}

func main() {
	parser := flags.NewParser(nil, flags.Default)

	unpackCmd := UnpackOptions{}
	_, _ = parser.AddCommand("unpack", "SH2 Unpacker", "Extracts files from SH2's game files", &unpackCmd)

	_, err := parser.Parse()
	handleFlagsError(err)
}
