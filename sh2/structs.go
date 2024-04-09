package sh2

import (
	"fmt"
)

type FileEntryType uint32

const (
	EntryTypeEOF        FileEntryType = 0x00
	EntryTypeBinaryFile FileEntryType = 0x03 // GX/GY/GZ files and such
	EntryTypeMergeFile  FileEntryType = 0x23 // mergefiles
	EntryTypeDataFile   FileEntryType = 0x50 // files within mergefiles

	MagicOffset = 0xFF800 // source: it came to me in a dream
)

type FilePathEntry struct {
	FileOffset uint32
	PathOffset uint32
}

func (e FilePathEntry) String() string {
	return fmt.Sprintf("{FileOffset:0x%X PathOffset:0x%X}", e.FileOffset, e.PathOffset)
}

// MergeFileEntry is used for both binary files and mergefiles.
// I'm calling it MergeFileEntry for simplicity.
// Unknown1 and Unknown2 seem to always be zero.
type MergeFileEntry struct {
	PathOffset uint32
	Unknown1   uint32
	Unknown2   uint32
}

func (e MergeFileEntry) String() string {
	return fmt.Sprintf("{PathOffset:0x%08X Unknown1:0x%08X Unknown2:0x%08X}", e.PathOffset, e.Unknown1, e.Unknown2)
}

// DataFileEntry describes where in a mergefile its contents can be found.
// EntryAddress technically points to a MergeFileEntry,
// but with 0x10 added to it for every consecutive new subdirectory.
// See DataMap.GetMergeFileEntryFromDataFileEntry for how this works in practice.
type DataFileEntry struct {
	EntryAddress uint32
	ChunkOffset  uint32
	ChunkLength  uint32
}

func (e DataFileEntry) String() string {
	return fmt.Sprintf("{EntryAddress:0x%08X ChunkOffset:0x%08X ChunkLength:0x%08X}", e.EntryAddress, e.ChunkOffset, e.ChunkLength)
}

type DataMap struct {
	FileToPathOffsets []FilePathEntry

	binaryFileOffsets map[uint32]MergeFileEntry
	mergeFileOffsets  map[uint32]MergeFileEntry
	dataFileOffsets   map[uint32]DataFileEntry

	filePaths map[uint32]string
}

// debugging functions
// func (d DataMap) GetBinaryFileEntries() []MergeFileEntry {
// 	return maps.Values(d.binaryFileOffsets)
// }
//
// func (d DataMap) GetMergeFileEntries() []MergeFileEntry {
// 	return maps.Values(d.mergeFileOffsets)
// }

func (d DataMap) GetBinaryFileEntry(rawAddress uint32) (MergeFileEntry, bool) {
	entry, ok := d.binaryFileOffsets[rawAddress-MagicOffset]
	if !ok {
		return MergeFileEntry{}, false
	}

	return entry, true
}

// GetMergeFileEntry takes a raw address and returns a MergeFileEntry.
func (d DataMap) GetMergeFileEntry(rawAddress uint32) (MergeFileEntry, bool) {
	entry, ok := d.mergeFileOffsets[rawAddress-MagicOffset]
	if !ok {
		return MergeFileEntry{}, false
	}

	return entry, true
}

// GetDataFileEntry takes a raw address and returns a DataFileEntry.
func (d DataMap) GetDataFileEntry(rawAddress uint32) (DataFileEntry, bool) {
	entry, ok := d.dataFileOffsets[rawAddress-MagicOffset]
	if !ok {
		return DataFileEntry{}, false
	}

	return entry, true
}

// GetMergeFileEntryFromDataFileEntry takes a DataFileEntry and returns a MergeFileEntry.
// This is done by taking the DataFileEntry's EntryAddress value and subtracting 0x10 in a loop until
// we have an address that matches a MergeFileEntry.
func (d DataMap) GetMergeFileEntryFromDataFileEntry(datEntry DataFileEntry, minOffset uint32) (MergeFileEntry, bool) {
	for addr := datEntry.EntryAddress; addr >= minOffset; addr -= 0x10 {
		entry, ok := d.mergeFileOffsets[addr-MagicOffset]
		if ok {
			return entry, true
		}
	}

	return MergeFileEntry{}, false
}

// GetFilePath takes a raw address and returns a file path string.
// This only works with FilePathEntry.FileOffset or MergeFileEntry.PathOffset addresses.
func (d DataMap) GetFilePath(rawAddress uint32) (string, bool) {
	path, ok := d.filePaths[rawAddress-MagicOffset]
	if !ok {
		return "", false
	}

	return path, true
}
