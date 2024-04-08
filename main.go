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

	log.Printf("Input: %s", opts.InFile)
	log.Printf("Output Folder: %s", opts.Pos.OutDir)

	inFileBytes, err := os.ReadFile(string(opts.InFile))
	if err != nil {
		log.Printf("ReadFile err: %v", err)
	}

	startOffset := utils.IndexOfSlice(inFileBytes, bin.PS2Padding())
	if startOffset < 0 {
		log.Fatalln("couldn't find data block.")
	}

	// skip special padding and null byte padding
	startOffset += 1024 + 24

	log.Printf("Data offset found at 0x%X", startOffset)

	dataMap, err := bin.ReadDataMap(string(opts.InFile), int64(startOffset))
	if err != nil {
		log.Printf("boop err: %v", err)
	}

	// iterate over the archive parts in the FTP list
	for _, ftp := range dataMap.FileToPathOffsets {
		arpEntry, arpOK := dataMap.GetArchivePartEntry(ftp.FileOffset)
		if !arpOK {
			continue
		}

		arpPath, ok := dataMap.GetFilePath(ftp.PathOffset)
		if !ok {
			log.Fatalf("Can't find file path for ARP at offset 0x%X", ftp.PathOffset)
		}

		arcEntry, ok := dataMap.GetArchiveFileEntryFromARPEntry(arpEntry, uint32(startOffset))
		if !ok {
			log.Fatalf("Can't find ARC entry for ARP %[2]s (%[3]s)", ftp.PathOffset, arpEntry, arpPath)
		}

		arcPath, ok := dataMap.GetFilePath(arcEntry.PathOffset)
		if !ok {
			log.Fatalf("Can't find file path for ARC %[2]s", ftp.PathOffset, arcEntry)
		}

		log.Printf("(%s) -> (%s)", arcPath, arpPath)
	}
}
