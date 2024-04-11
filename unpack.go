package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sh2unpack/sh2"
	"sh2unpack/utils"
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

	gameVersion, ok := sh2.VersionMap[shaString]
	if !ok {
		return fmt.Errorf("Not a supported file or gameVersion of the game: %s", inFilePath)
	}

	fmt.Printf("Version detected: %s, %s\n", gameVersion.FileName, gameVersion.Description)

	dataMap, err := sh2.ReadDataMap(inFile, gameVersion, opts.Debug)
	if err != nil {
		return fmt.Errorf("Couldn't read data map: %v", err)
	}

	// explicitly close the input file here, it's no longer needed
	_ = inFile.Close()

	if opts.Debug {
		guessedOffset := dataMap.GuessOffset()
		fmt.Printf("guessed offset: 0x%X\n", guessedOffset)
		fmt.Printf("actual offset:  0x%X\n", gameVersion.MagicOffset)
	}

	mergeFileMap := map[string]*os.File{}
	defer func() {
		for _, f := range mergeFileMap {
			_ = f.Close()
		}
	}()

	if opts.DryRun {
		fmt.Println("Doing a dry run.")
	}

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

		mgfEntry, ok := dataMap.GetMergeFileEntryFromDataFileEntry(datEntry, gameVersion.DataOffset)
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

		destinationPath := filepath.Join(outDirPath, datPath)
		destinationDir := filepath.Dir(destinationPath)
		mgfBase := filepath.Base(mgfPath)

		if !opts.DryRun {
			err = os.MkdirAll(destinationDir, 0700)
			if err != nil {
				return fmt.Errorf("Can't create destination dir %s: %v", destinationDir, err)
			}

			f, err := os.Create(destinationPath)
			if err != nil {
				return fmt.Errorf("Can't create destination file %s: %v", destinationPath, err)
			}

			err = utils.CopyPartOfFileToFile(f, mergeFile, int64(datEntry.ChunkOffset), int64(datEntry.ChunkLength))
			if err != nil {
				return fmt.Errorf("Can't copy chunk from %s to %s: %v", mgfBase, destinationPath, err)
			}
			_ = f.Close()

			if opts.Debug {
				fmt.Printf("Extracted %d bytes from %s to %s\n", datEntry.ChunkLength, mgfBase, destinationPath)
			}
		}

		numExtractedFiles++

	}

	fmt.Printf("Extracted %d files.\n", numExtractedFiles)

	return nil
}
