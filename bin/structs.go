package bin

import "fmt"

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
	return fmt.Sprintf("{EntryAddress:0x%X ChunkOffset:0x%X ChunkLength:0x%X}", e.EntryAddress, e.ChunkOffset, e.ChunkLength)
}

type DataMap struct {
	FileToPathPointers map[uint32]uint32

	BinaryFilePointers      map[uint32]Table2FileEntry
	ArchiveFilePointers     map[uint32]Table2FileEntry
	ArchiveDeepFilePointers map[uint32]Table2ChunkEntry

	FilePaths map[uint32]string
}
