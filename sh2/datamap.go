package sh2

import (
	"errors"
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

func ReadDataMap(f *os.File, startOffset int64) (*DataMap, error) {
	pos, err := f.Seek(startOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	dataMap := DataMap{
		FileToPathOffsets: []FilePathEntry{},
		binaryFileOffsets: map[uint32]MergeFileEntry{},
		mergeFileOffsets:  map[uint32]MergeFileEntry{},
		dataFileOffsets:   map[uint32]DataFileEntry{},
		filePaths:         map[uint32]string{},
	}

	for {
		var entry FilePathEntry
		err := utils.ReadStruct(f, &entry)
		if err != nil {
			return nil, err
		}

		if entry.FileOffset == 0 && entry.PathOffset == 0 {
			pos, _ = f.Seek(-8, io.SeekCurrent)
			break
		}

		dataMap.FileToPathOffsets = append(dataMap.FileToPathOffsets, entry)
	}

	// fmt.Printf("total file-path entries: %d\n", len(dataMap.FileToPathOffsets))

	// skip to the next table
	pos, err = f.Seek(80, io.SeekCurrent) // TODO: check if this offset is different in other game versions
	if err != nil {
		return nil, err
	}

	for {
		pos, _ = f.Seek(0, io.SeekCurrent)
		posUint := uint32(pos)

		var entryType FileEntryType
		err := utils.ReadStruct(f, &entryType)
		if err != nil {
			return nil, err
		}

		switch entryType {
		case EntryTypeBinaryFile:
			var entry MergeFileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.binaryFileOffsets[posUint] = entry
		case EntryTypeMergeFile:
			var entry MergeFileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.mergeFileOffsets[posUint] = entry
		case EntryTypeDataFile:
			var entry DataFileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			dataMap.dataFileOffsets[posUint] = entry
		}

		if entryType == EntryTypeEOF {
			pos, _ = f.Seek(-4, io.SeekCurrent)
			break
		}
	}

	// fmt.Printf("total binary files: %d\n", len(dataMap.binaryFileOffsets))
	// fmt.Printf("total mergefiles: %d\n", len(dataMap.mergeFileOffsets))
	// fmt.Printf("total data files: %d\n", len(dataMap.dataFileOffsets))

	// skip to the next table
	pos, err = f.Seek(48, io.SeekCurrent) // TODO: check if this offset is different in other game versions
	if err != nil {
		return nil, err
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

	// fmt.Printf("total file paths: %d\n", len(dataMap.filePaths))

	return &dataMap, nil
}
