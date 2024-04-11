package sh2

import (
	"errors"
	"fmt"
	"io"
	"os"

	"sh2unpack/utils"
)

/*
	[NTSC 2.01] SLUS_202.28
	Table 1 Offset: 0x2CCF00 (2936576)
	Table 2 Offset: 0x2D4900 (2967808)
	Table 3 Offset: 0x2E4680 (3032704)
	Data Table End: 0x2FFCE0 (3144928)
	Data Table Len: 0x032DE0 ( 208352)
*/

func skipToNextTable(f *os.File, maxSteps int, debug bool) error {
	// skip to the next table by advancing in 8-byte steps until non-null bytes are found
	for i := 0; i < maxSteps; i++ {
		var sentinel uint64
		err := utils.ReadStructLE(f, &sentinel)
		if err != nil {
			return err
		}

		if sentinel != 0 {
			if debug {
				fmt.Printf("skipped %d bytes\n", i*8)
			}

			_, _ = f.Seek(-8, io.SeekCurrent)
			break
		}
	}

	return nil
}

func ReadDataMap(f *os.File, gv gameVersion, debug bool) (*DataMap, error) {
	pos, err := f.Seek(int64(gv.DataOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	dataMap := DataMap{
		FileToPathOffsets: []FilePathEntry{},
		binaryFileOffsets: map[uint32]MergeFileEntry{},
		mergeFileOffsets:  map[uint32]MergeFileEntry{},
		dataFileOffsets:   map[uint32]DataFileEntry{},
		filePaths:         map[uint32]string{},
		magicOffset:       gv.MagicOffset,
	}

	if debug {
		fmt.Printf("start offset: 0x%X\n", pos)
	}

	for {
		var entry FilePathEntry
		err := utils.ReadStructLE(f, &entry)
		if err != nil {
			return nil, err
		}

		if entry.FileOffset == 0 && entry.PathOffset == 0 {
			_, _ = f.Seek(-8, io.SeekCurrent)
			break
		}

		dataMap.FileToPathOffsets = append(dataMap.FileToPathOffsets, entry)
	}

	if debug {
		fmt.Printf("total file-path entries: %d\n", len(dataMap.FileToPathOffsets))
	}

	err = skipToNextTable(f, 32, debug)
	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf("pos before file entry table: 0x%X\n", utils.CurrentPos(f))
	}

	for {
		posUint := uint32(utils.CurrentPos(f))

		var entryType FileEntryType
		err := utils.ReadStructLE(f, &entryType)
		if err != nil {
			return nil, err
		}

		switch entryType {
		case EntryTypeBinaryFile:
			var entry MergeFileEntry
			err := utils.ReadStructLE(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.binaryFileOffsets[posUint] = entry
		case EntryTypeMergeFile:
			var entry MergeFileEntry
			err := utils.ReadStructLE(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.mergeFileOffsets[posUint] = entry
		case EntryTypeDataFile:
			var entry DataFileEntry
			err := utils.ReadStructLE(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.dataFileOffsets[posUint] = entry
		}

		if entryType == EntryTypeEOF {
			_, _ = f.Seek(-4, io.SeekCurrent)
			break
		}
	}

	if debug {
		fmt.Printf("total binary files: %d\n", len(dataMap.binaryFileOffsets))
		fmt.Printf("total mergefiles: %d\n", len(dataMap.mergeFileOffsets))
		fmt.Printf("total data files: %d\n", len(dataMap.dataFileOffsets))
	}

	// skip to the next table
	err = skipToNextTable(f, 32, debug)
	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf("pos before path table: 0x%X\n", utils.CurrentPos(f))
	}

	for {
		pathOffset, pathEntry, err := utils.ReadNullTerminatedString(f)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidASCIIChar) || errors.Is(err, io.EOF) {
				break // reached the end of the path table
			}

			return nil, err
		}

		dataMap.filePaths[uint32(pathOffset)] = pathEntry
	}

	if debug {
		fmt.Printf("total file paths: %d\n", len(dataMap.filePaths))
	}

	return &dataMap, nil
}
