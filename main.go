package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"sh2unpack/constants"
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

	// print version
	fmt.Printf("sh2unpack %s [%s]\n", constants.GitVersion, constants.GitCommitShort)

	_, err := parser.Parse()
	handleFlagsError(err)
}
