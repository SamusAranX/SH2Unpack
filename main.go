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

	dataMap, err := bin.ReadDataMap(string(opts.InFile))
	if err != nil {
		log.Printf("boop err: %v", err)
	}

	// minFTPPointer := slices.Min(maps.Values(dataMap.FileToPathPointers))
	// minPathKey := slices.Min(maps.Keys(dataMap.FilePaths))
	// log.Printf("min ftp pointer: %d", minFTPPointer)
	// log.Printf("min path pointer: %d", minPathKey)
	// return

	for filePointer, pathPointer := range dataMap.FileToPathPointers {
		mangledPointer := filePointer - bin.MagicOffset

		filePath, pathOK := dataMap.FilePaths[pathPointer-bin.PathPointerOffset]
		if !pathOK {
			log.Printf("%X - %X = %X", filePointer, bin.PathPointerOffset, mangledPointer)
			log.Fatalln("file path not found!")
		}

		binEntry, binOK := dataMap.BinaryFilePointers[mangledPointer]
		arcEntry, arcOK := dataMap.ArchiveFilePointers[mangledPointer]
		dpfEntry, dpfOK := dataMap.ArchiveDeepFilePointers[mangledPointer]
		if binOK {
			log.Printf("[BIN|   |   ] %s (%s)", binEntry, filePath)
		} else if arcOK {
			log.Printf("[   |ARC|   ] %s (%s)", arcEntry, filePath)
		} else if dpfOK {
			log.Printf("[   |   |DPF] %s (%s)", dpfEntry, filePath)
		} else {
			log.Fatalln("no file found!")
		}
	}
}
