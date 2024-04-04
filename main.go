package main

import (
	"errors"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"sh2unpack/bin"
)

func handleFlagsError(err error) {
	if err != nil {
		var e *flags.Error
		if errors.As(err, &e) {
			// log.Printf("type: %d", e.Type)
			switch {
			case errors.Is(e, flags.ErrHelp):
				os.Exit(0)
			case errors.Is(e, flags.ErrRequired):
				os.Exit(1)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func main() {
	var opts UnpackOptions
	_, err := flags.Parse(&opts)
	handleFlagsError(err)

	// log.Printf("in:  %s", opts.InFile)
	// log.Printf("out: %s", opts.Pos.OutDir)

	err = bin.Boop(string(opts.InFile), string(opts.Pos.OutDir))
	if err != nil {
		log.Printf("boop err: %v", err)
	}
}
