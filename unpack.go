package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sh2unpack/sh2"
	"sh2unpack/utils"
)

type gameVersion struct {
	DataOffset  uint32
	FileName    string
	Description string
}

var (
	// map of game binary SHA1 hash -> game version
	versionMap = map[string]gameVersion{
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
	}
)

func (opts *UnpackOptions) Execute(args []string) error {
	inFilePath := string(opts.InFile)
	outDirPath := string(opts.Pos.OutDir)

	fmt.Printf("Input: %s\n", inFilePath)
	fmt.Printf("Output Folder: %s\n", outDirPath)

	inFile, err := os.Open(inFilePath)
	if err != nil {
		return fmt.Errorf("Can't open file: %v", err)
	}

	shaString, err := utils.HashFileSHA1(inFile)
	if err != nil {
		return fmt.Errorf("Can't hash input file: %v", err)
	}

	version, ok := versionMap[shaString]
	if !ok {
		return fmt.Errorf("Not a supported file or version of the game: %s", inFilePath)
	}

	fmt.Printf("Version detected: %s, %s\n", version.FileName, version.Description)

	dataMap, err := sh2.ReadDataMap(inFile, int64(version.DataOffset))
	if err != nil {
		return fmt.Errorf("Couldn't read data map: %v", err)
	}

	// explicitly close the input file here, it's no longer needed
	_ = inFile.Close()

	mergeFileMap := map[string]*os.File{}
	defer func() {
		for _, f := range mergeFileMap {
			_ = f.Close()
		}
	}()

	numExtractedFiles := 0

	// iterate over the data files in the FTP list
	for _, ftp := range dataMap.FileToPathOffsets {
		datEntry, ok := dataMap.GetDataFileEntry(ftp.FileOffset)
		if !ok {
			continue
		}

		datPath, ok := dataMap.GetFilePath(ftp.PathOffset)
		if !ok {
			return fmt.Errorf("Can't find file path for data file at offset 0x%X", ftp.PathOffset)
		}

		mgfEntry, ok := dataMap.GetMergeFileEntryFromDataFileEntry(datEntry, version.DataOffset)
		if !ok {
			return fmt.Errorf("Can't find mergefile entry for data file %[2]s (%[3]s)", ftp.PathOffset, datEntry, datPath)
		}

		mgfPath, ok := dataMap.GetFilePath(mgfEntry.PathOffset)
		if !ok {
			return fmt.Errorf("Can't find file path for mergefile %[2]s", ftp.PathOffset, mgfEntry)
		}

		mergeFile, ok := mergeFileMap[mgfPath]
		if !ok {
			actualMGFPath := filepath.Join(filepath.Dir(inFilePath), strings.ToUpper(mgfPath))

			f, err := os.Open(actualMGFPath)
			if err != nil {
				return fmt.Errorf("Can't open mergefile: %v", err)
			}
			mergeFileMap[mgfPath] = f
			mergeFile = f
		}

		destinationPath := filepath.Join(outDirPath, strings.ToUpper(datPath))
		destinationDir := filepath.Dir(destinationPath)

		err = os.MkdirAll(destinationDir, 0700)
		if err != nil {
			return fmt.Errorf("Can't create destination dir %s: %v", destinationDir, err)
		}

		f, err := os.Create(destinationPath)
		if err != nil {
			return fmt.Errorf("Can't create destination file %s: %v", destinationPath, err)
		}

		mgfBase := filepath.Base(mgfPath)

		err = utils.CopyPartOfFileToFile(f, mergeFile, int64(datEntry.ChunkOffset), int64(datEntry.ChunkLength))
		if err != nil {
			return fmt.Errorf("Can't copy chunk from %s to %s: %v", mgfBase, destinationPath, err)
		}
		_ = f.Close()
		numExtractedFiles++

		if opts.Debug {
			fmt.Printf("Extracted %d bytes from %s to %s\n", datEntry.ChunkLength, mgfBase, destinationPath)
		}
	}

	fmt.Printf("Extracted %d files.\n", numExtractedFiles)

	return nil
}
