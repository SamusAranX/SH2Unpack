package bin

import (
	"fmt"
	"io"
	"log"
	"os"

	"sh2unpack/utils"
)

const (
	tableStart  = 0x2CCF00 // this is hardcoded for NTSC 2.01
	magicOffset = 0xFF800  // source: it came to me in a dream
)

/*
	Table 1 Offset: 0x2CCF00 (2936576)
	Table 2 Offset: 0x2D4900 (2967808)
	Table 3 Offset: 0x2E4680 (3032704)
	Data Table End: 0x2FFCF3 (3144947)
*/

func Boop(inFile, outDir string) error {
	f, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}

	defer f.Close()

	pos, err := f.Seek(tableStart, io.SeekStart)
	if err != nil {
		return err
	}

	var (
		table1Entries = map[uint32]uint32{} // map of filePointer -> pathPointer

		table2PhysFileEntries = map[uint32]Table2FileEntry{}
		table2VirtFileEntries = map[uint32]Table2FileEntry{}
		table2ChunkEntries    = map[uint32]Table2ChunkEntry{}
		table3FilePaths       = map[uint32]string{} // map of offset in binary -> file path
	)

	for {
		var entry Table1Entry
		err := utils.ReadStruct(f, &entry)
		if err != nil {
			return err
		}

		if entry.FilePointer == 0 && entry.PathPointer == 0 {
			break
		}

		table1Entries[entry.FilePointer] = entry.PathPointer
	}

	table1EntriesCount := len(table1Entries)
	log.Printf("Read %d entries into table 1", table1EntriesCount)

	// skip to the next table
	pos, err = f.Seek(72, io.SeekCurrent)
	if err != nil {
		return err
	}

	for {
		pos, _ = f.Seek(0, io.SeekCurrent)
		posUint := uint32(pos)

		var entryType Table2EntryType
		err := utils.ReadStruct(f, &entryType)
		if err != nil {
			return err
		}

		switch entryType {
		case EntryTypePhysicalFile:
			var entry Table2FileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return err
			}

			table2PhysFileEntries[posUint] = entry
		case EntryTypeVirtualFile:
			var entry Table2FileEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return err
			}

			table2VirtFileEntries[posUint] = entry
		case EntryTypeVirtualChunk:
			var entry Table2ChunkEntry
			err := utils.ReadStruct(f, &entry)
			if err != nil {
				return err
			}

			table2ChunkEntries[posUint] = entry
		}

		if entryType == EntryTypeEOF {
			break
		}

		// if entry.FileOffset == 0x49F800 {
		// 	// forestPathOffset := 0x2FD490
		// 	log.Printf("forest.sfc: %+v", entry)
		// 	log.Printf("  EntryType: %d", entry.EntryType)
		// 	log.Printf("  PathOffset: 0x%X (%[1]d)", entry.PathOffset)
		// 	log.Printf("  FileOffset: 0x%X (%[1]d)", entry.FileOffset)
		// 	log.Printf("  FileLength: 0x%X (%[1]d)", entry.FileLength)
		// }
	}

	table2PhysFileEntriesCount := len(table2PhysFileEntries)
	table2VirtFileEntriesCount := len(table2VirtFileEntries)
	table2ChunkEntriesCount := len(table2ChunkEntries)

	log.Printf("Read %d physical file entries from table 2", table2PhysFileEntriesCount)
	log.Printf("Read %d virtual file entries from table 2", table2VirtFileEntriesCount)
	log.Printf("Read %d chunk entries from table 2", table2ChunkEntriesCount)
	log.Printf("Read %d entries from table 2 in total", table2PhysFileEntriesCount+table2VirtFileEntriesCount+table2ChunkEntriesCount)

	// skip to the next table
	pos, err = f.Seek(44, io.SeekCurrent)
	if err != nil {
		return err
	}

	for {
		pathOffset, pathEntry, err := utils.ReadNullTerminatedString(f)
		if err != nil {
			if err == utils.ErrInvalidASCIIChar {
				break // reached the end of the path table
			}

			return err
		}

		table3FilePaths[uint32(pathOffset)] = pathEntry
	}

	table3FilePathsCount := len(table3FilePaths)
	log.Printf("Read %d file paths from table 3", table3FilePathsCount)

	// check where we are
	pos, _ = f.Seek(0, io.SeekCurrent)

	return nil
}
