package bin

import (
	"fmt"
)

type Table2EntryType uint32

const (
	EntryTypeEOF          = 0
	EntryTypePhysicalFile = 0x3  // GX/GY/GZ files and such
	EntryTypeVirtualFile  = 0x23 // MGF files
	EntryTypeVirtualChunk = 0x50 // files within MGF files

	// this is hardcoded for NTSC 2.01
	//TableStartOffset = 0x2CCF00

	MagicOffset = 0xFF800 // source: it came to me in a dream
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
	return fmt.Sprintf("{PathOffset:0x%08X FileOffset:0x%08X FileLength:0x%08X}", e.PathOffset, e.FileOffset, e.FileLength)
}

// used for chunks (or "sub files")
type Table2ChunkEntry struct {
	EntryAddress uint32
	ChunkOffset  uint32
	ChunkLength  uint32
}

func (e Table2ChunkEntry) String() string {
	return fmt.Sprintf("{EntryAddress:0x%08X ChunkOffset:0x%08X ChunkLength:0x%08X}", e.EntryAddress, e.ChunkOffset, e.ChunkLength)
}

type DataMap struct {
	FileToPathPointers []Table1Entry

	BinaryFilePointers  map[uint32]Table2FileEntry
	ArchiveFilePointers map[uint32]Table2FileEntry
	ArchivePartPointers map[uint32]Table2ChunkEntry

	FilePaths map[uint32]string
}

//	func (d DataMap) GetBinaryFileEntry(rawAddress uint32) Table2FileEntry {
//		entry, ok := d.BinaryFilePointers[rawAddress-BINAddressOffset]
//		if !ok {
//			log.Fatalf("invalid binary file pointer 0x%X (%[1]d)", rawAddress)
//		}
//
//		return entry
//	}

// GetArchiveFileEntry takes a raw address and returns a struct representing an MGF file
func (d DataMap) GetArchiveFileEntry(rawAddress uint32) (Table2FileEntry, bool) {
	entry, ok := d.ArchiveFilePointers[rawAddress-MagicOffset]
	//entry, ok := d.ArchiveFilePointers[rawAddress]
	if !ok {
		// log.Printf("invalid archive file pointer 0x%X (%[1]d)", rawAddress)
		return Table2FileEntry{}, false
	}

	return entry, true
}

// GetArchiveFileEntryFromARPEntry takes an ARP entry and returns a struct representing an MGF file
func (d DataMap) GetArchiveFileEntryFromARPEntry(arpEntry Table2ChunkEntry) (Table2FileEntry, bool) {
	entry, ok := d.ArchiveFilePointers[arpEntry.EntryAddress-MagicOffset]
	//entry, ok := d.ArchiveFilePointers[rawAddress]
	if !ok {
		// log.Printf("invalid archive file pointer 0x%X (%[1]d)", rawAddress)
		return Table2FileEntry{}, false
	}

	return entry, true
}

// GetArchivePartEntry takes a raw address and returns a struct representing a file stored inside of an MGF file
func (d DataMap) GetArchivePartEntry(rawAddress uint32) (Table2ChunkEntry, bool) {
	entry, ok := d.ArchivePartPointers[rawAddress-MagicOffset]
	if !ok {
		// log.Printf("invalid archive part pointer 0x%X (%[1]d)", rawAddress)
		return Table2ChunkEntry{}, false
	}

	return entry, true
}

// GetFilePath takes a raw address and returns a file path string
func (d DataMap) GetFilePath(rawAddress uint32) (string, bool) {
	path, ok := d.FilePaths[rawAddress-MagicOffset]
	if !ok {
		return "", false
		// log.Fatalf("invalid file path pointer 0x%X (%[1]d)", rawAddress)
	}

	return path, true
}
