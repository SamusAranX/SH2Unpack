package main

import (
	"errors"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	var opts UnpackOptions
	_, err := flags.Parse(&opts)

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

	log.Println(opts)
}
