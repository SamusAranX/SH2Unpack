package main

import (
	"errors"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
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

func nop(a ...any) {}

func main() {
	var opts UnpackOptions
	_, err := flags.Parse(&opts)
	handleFlagsError(err)

	log.Println(opts.InFile)
	// log.Printf("out: %s", opts.Pos.OutDir)

	dataMap, err := bin.ReadDataMap(string(opts.InFile), 0x2CCF00)
	if err != nil {
		log.Printf("boop err: %v", err)
	}

	for _, ftp := range dataMap.FileToPathPointers {
		path, pathOK := dataMap.GetFilePath(ftp.PathPointer)
		if !pathOK {
			log.Fatalf("Can't find path for offset 0x%X!", ftp.PathPointer)
		}

		arcEntry, arcOK := dataMap.GetArchiveFileEntry(ftp.FilePointer)
		arpEntry, arpOK := dataMap.GetArchivePartEntry(ftp.FilePointer)
		nop(arcEntry, arpEntry)

		if arcOK {
			log.Printf("[ARC]   0x%08X %s (%s)", ftp.FilePointer, arcEntry, path)
		} else if arpOK {
			log.Printf("  {arp} 0x%08X %[2]s (%[3]s)", ftp.FilePointer, arpEntry, path)
		}
	}
}
