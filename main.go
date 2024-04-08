package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"path/filepath"
	"sh2unpack/bin"
	"sh2unpack/utils"
	"strings"
)

type gameVersion struct {
	DataOffset  uint32
	FileName    string
	Description string
}

var (
	versionMap = map[string]gameVersion{
		// Demos/Prototypes
		"B9CB2E895FC83CD4452DC9A818BF3CA26394ADBE": {
			DataOffset:  0x2B3120,
			FileName:    "SLPM_610.09",
			Description: "PAL v1.00 (Trial Version)",
		},
		"50C664C525736619215654186446A5D6B211FB31": {
			DataOffset:  0x45C200,
			FileName:    "SLPM_123.45",
			Description: "NTSC v0.30 (E3 Demo)",
		},
		"888EFF71606FF4C1C610E30111B3CA5DA647EDCC": {
			DataOffset:  0x29CD00,
			FileName:    "SLUS_202.28",
			Description: "NTSC v0.10 (2001-07-13 prototype)",
		},

		// NTSC
		"3A27DEDDFA81CF30F46F0742C3523230CAC75D9A": {
			DataOffset:  0x2CCF00,
			FileName:    "SLUS_202.28",
			Description: "NTSC v2.01 (Greatest Hits)",
		},

		// PAL
		"8BC367E1B9E7AA5CC5D5FA32048ED97F3FADE728": {
			DataOffset:  0x2BD390,
			FileName:    "SLES_503.82",
			Description: "PAL v1.02 (Base Game)",
		},
		"2C5A7AFBA3A5B4507CCB828811C8ADD9E5D0E961": {
			DataOffset:  0x2CD980,
			FileName:    "SLES_511.56",
			Description: "PAL v1.10 (Director's Cut)",
		},
	}
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

	inFile := string(opts.InFile)
	outDir := string(opts.Pos.OutDir)

	log.Printf("Input: %s", inFile)
	log.Printf("Output Folder: %s", outDir)

	f, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("Can't open file: %v", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("Can't hash input file: %v", err)
	}

	shaString := fmt.Sprintf("%X", h.Sum(nil))

	version, ok := versionMap[shaString]
	if !ok {
		log.Fatalf("Not a supported file or version of the game: %s", inFile)
	}

	log.Printf("Version detected: %s, %s", version.FileName, version.Description)

	dataMap, err := bin.ReadDataMap(f, int64(version.DataOffset))
	if err != nil {
		log.Printf("boop err: %v", err)
	}

	archiveMap := map[string]*os.File{}
	defer func() {
		for _, f := range archiveMap {
			_ = f.Close()
		}
	}()

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

		arcEntry, ok := dataMap.GetArchiveFileEntryFromARPEntry(arpEntry, uint32(version.DataOffset))
		if !ok {
			log.Fatalf("Can't find ARC entry for ARP %[2]s (%[3]s)", ftp.PathOffset, arpEntry, arpPath)
		}

		arcPath, ok := dataMap.GetFilePath(arcEntry.PathOffset)
		if !ok {
			log.Fatalf("Can't find file path for ARC %[2]s", ftp.PathOffset, arcEntry)
		}

		arcFile, ok := archiveMap[arcPath]
		if !ok {
			actualArcPath := filepath.Join(filepath.Dir(inFile), strings.ToUpper(arcPath))

			f, err := os.Open(actualArcPath)
			if err != nil {
				log.Fatalf("Can't open file: %v", err)
			}
			archiveMap[arcPath] = f
			arcFile = f
		}

		actualDestinationPath := filepath.Join(outDir, strings.ToUpper(arpPath))

		err = os.MkdirAll(filepath.Dir(actualDestinationPath), 0700)
		if err != nil {
			log.Fatalf("Can't create destination dir: %v", err)
		}

		nop(arcFile)

		f, err := os.Create(actualDestinationPath)
		if err != nil {
			log.Fatalf("Can't create destination file: %v", err)
		}

		arcBase := filepath.Base(arcPath)
		dstBase := filepath.Base(actualDestinationPath)

		err = utils.CopyPartOfFileToFile(f, arcFile, int64(arpEntry.ChunkOffset), int64(arpEntry.ChunkLength))
		if err != nil {
			log.Fatalf("Can't copy chunk from %s to %s: %v", arcBase, dstBase, err)
		}
		_ = f.Close()

		log.Printf("Extracted %d bytes to %s", arpEntry.ChunkLength, actualDestinationPath)
	}

	log.Println("breakpoint here")
}
