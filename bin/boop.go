package bin

import (
	"errors"
	"fmt"
	"io"
	"log"
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

func ReadDataMap(inFile string, startOffset int64) (*DataMap, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %v", err)
	}

	defer f.Close()

	pos, err := f.Seek(startOffset, io.SeekStart)
	//pos, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var (
		table1Entries []Table1Entry

		table2PhysFileEntries = map[uint32]Table2FileEntry{}  // map of offset in binary -> entry
		table2VirtFileEntries = map[uint32]Table2FileEntry{}  // map of offset in binary -> entry
		table2ChunkEntries    = map[uint32]Table2ChunkEntry{} // map of offset in binary -> entry
		table3FilePaths       = map[uint32]string{}           // map of offset in binary -> file path
	)

	for {
		var entry Table1Entry
		err := utils.ReadStruct(f, &entry)
		if err != nil {
			return nil, err
		}

		if entry.FileOffset == 0 && entry.PathOffset == 0 {
			pos, _ = f.Seek(-8, io.SeekCurrent)
			break
		}

		table1Entries = append(table1Entries, entry)
	}

	table1EntriesCount := len(table1Entries)
	log.Printf("[File/Path Table] total entries: %d", table1EntriesCount)

	// skip to the next table
	pos, err = f.Seek(80, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	for {
		pos, _ = f.Seek(0, io.SeekCurrent)
		posUint := uint32(pos)

		var entryType Table2EntryType
		err := utils.ReadStruct(f, &entryType)
		if err != nil {
			return nil, err
		}

		switch entryType {
		case EntryTypePhysicalFile:
			var entry Table2FileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			table2PhysFileEntries[posUint] = entry
		case EntryTypeVirtualFile:
			var entry Table2FileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			table2VirtFileEntries[posUint] = entry
		case EntryTypeVirtualChunk:
			var entry Table2ChunkEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return nil, err
			}

			table2ChunkEntries[posUint] = entry
		}

		if entryType == EntryTypeEOF {
			pos, _ = f.Seek(-4, io.SeekCurrent)
			break
		}
	}

	table2PhysFileEntriesCount := len(table2PhysFileEntries)
	table2VirtFileEntriesCount := len(table2VirtFileEntries)
	table2ChunkEntriesCount := len(table2ChunkEntries)

	log.Printf("[File Table] phys files: %d", table2PhysFileEntriesCount)
	log.Printf("[File Table] virt files: %d", table2VirtFileEntriesCount)
	log.Printf("[File Table] chunk files: %d", table2ChunkEntriesCount)

	// skip to the next table
	pos, err = f.Seek(48, io.SeekCurrent)
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

		table3FilePaths[uint32(pathOffset)] = pathEntry
	}

	table3FilePathsCount := len(table3FilePaths)
	log.Printf("[Path Table] file paths: %d", table3FilePathsCount)

	return &DataMap{
		FileToPathOffsets:  table1Entries,
		BinaryFileOffsets:  table2PhysFileEntries,
		ArchiveFileOffsets: table2VirtFileEntries,
		ArchivePartOffsets: table2ChunkEntries,
		FilePaths:          table3FilePaths,
	}, nil
}
