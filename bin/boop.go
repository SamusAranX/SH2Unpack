package bin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	paddingLength = 256
	tableStart    = 0x2CCF00
	weirdConstant = 0xFF800

	// table1EntryLength = 8
	// table2EntryLength = 16
)

type Table2EntryType uint32

const (
	EntryTypeEOF          = 0
	EntryTypePhysicalFile = 0x3
	EntryTypeVirtualFile  = 0x23
	EntryTypeVirtualChunk = 0x50
)

type Table1Entry struct {
	FilePointer uint32
	PathPointer uint32
}

func (e Table1Entry) String() string {
	return fmt.Sprintf("{FilePointer:0x%X PathPointer:0x%X}", e.FilePointer, e.PathPointer)
}

// used for physical and virtual files
type Table2FileEntry struct {
	PathOffset uint32
	FileOffset uint32
	FileLength uint32
}

func (e Table2FileEntry) String() string {
	return fmt.Sprintf("{PathOffset:0x%X FileOffset:0x%X FileLength:0x%X}", e.PathOffset, e.FileOffset, e.FileLength)
}

// used for chunks (or "sub files")
type Table2ChunkEntry struct {
	EntryAddress uint32
	ChunkOffset  uint32
	ChunkLength  uint32
}

func (e Table2ChunkEntry) String() string {
	return fmt.Sprintf("{PathOffset:0x%X FileOffset:0x%X FileLength:0x%X}", e.EntryAddress, e.ChunkOffset, e.ChunkLength)
}

func ps2Padding() [paddingLength]byte {
	// 1+1+2+4+8+16+32+64+128
	var padding [paddingLength]byte
	for i := 0; i < 1; i++ {
		padding[i] = 0
	}
	for i := 1; i < 2; i++ {
		padding[i] = 1
	}
	for i := 2; i < 2+2; i++ {
		padding[i] = 2
	}
	for i := 4; i < 4+4; i++ {
		padding[i] = 3
	}
	for i := 8; i < 8+8; i++ {
		padding[i] = 4
	}
	for i := 16; i < 16+16; i++ {
		padding[i] = 5
	}
	for i := 32; i < 32+32; i++ {
		padding[i] = 6
	}
	for i := 64; i < 64+64; i++ {
		padding[i] = 7
	}
	for i := 128; i < 128+128; i++ {
		padding[i] = 8
	}

	return padding
}

func readStruct[T interface{}](r io.Reader, t *T) error {
	err := binary.Read(r, binary.LittleEndian, t)
	if err != nil {
		return err
	}

	return nil
}

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

		table2PhysFileEntries []Table2FileEntry
		table2VirtFileEntries []Table2FileEntry
		table2ChunkEntries    []Table2ChunkEntry
		table3FilePaths       []string
	)

	log.Printf("Table 1 Offset: 0x%X (%[1]d)", pos)

	for {
		var entry Table1Entry
		err := readStruct(f, &entry)
		if err != nil {
			return err
		}

		if entry.FilePointer == 0 && entry.PathPointer == 0 {
			break
		}

		_, duplicateKey := table1Entries[entry.FilePointer]
		if duplicateKey {
			return errors.New("duplicate file pointer in table 1")
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

	log.Printf("Table 2 Offset: 0x%X (%[1]d)", pos)

	for {
		var entryType Table2EntryType
		err := readStruct(f, &entryType)
		if err != nil {
			return err
		}

		switch entryType {
		case EntryTypePhysicalFile:
			var entry Table2FileEntry
			err := readStruct(f, &entry)
			if err != nil {
				return err
			}

			table2PhysFileEntries = append(table2PhysFileEntries, entry)
		case EntryTypeVirtualFile:
			var entry Table2FileEntry
			err := readStruct(f, &entry)
			if err != nil {
				return err
			}

			table2VirtFileEntries = append(table2VirtFileEntries, entry)
		case EntryTypeVirtualChunk:
			var entry Table2ChunkEntry
			err := readStruct(f, &entry)
			if err != nil {
				return err
			}

			table2ChunkEntries = append(table2ChunkEntries, entry)
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

	log.Printf("Table 3 Offset: 0x%X (%[1]d)", pos)

	// for _, entry := range table2PhysFileEntries {
	// 	pos, _ = f.Seek(int64(entry.PathOffset-weirdConstant), io.SeekStart)
	// 	pathStr, _ := readNullTerminatedString(f)
	// 	// log.Printf("→ offset 0x%X yielded string \"%s\"", pos, pathStr)
	// 	log.Printf("phys file %s: %s", entry, pathStr)
	// }
	//
	// log.Println()
	//
	// for _, entry := range table2VirtFileEntries {
	// 	pos, _ = f.Seek(int64(entry.PathOffset-weirdConstant), io.SeekStart)
	// 	pathStr, _ := readNullTerminatedString(f)
	// 	// log.Printf("→ offset 0x%X yielded string \"%s\"", pos, pathStr)
	// 	log.Printf("virt file %s: %s", entry, pathStr)
	// }

	for {
		pathEntry, err := readNullTerminatedString(f)
		if err != nil {
			if err == ErrInvalidASCIIChar {
				break
			}

			return err
		}

		table3FilePaths = append(table3FilePaths, pathEntry)
	}

	table3FilePathsCount := len(table3FilePaths)
	log.Printf("Read %d file paths from table 3", table3FilePathsCount)

	// check where we are
	pos, _ = f.Seek(0, io.SeekCurrent)

	log.Printf("Ended up at 0x%X (%[1]d)", pos)

	return nil
}
