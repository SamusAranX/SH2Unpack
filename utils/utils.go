package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"os"

	"golang.org/x/exp/slices"
)

var (
	ErrInvalidASCIIChar  = errors.New("invalid ASCII char")
	ErrMaxLengthExceeded = errors.New("max search length exceeded")
)

// ReadStructLE reads arbitrary types and structs from io.Readers (little endian).
func ReadStructLE[T interface{}](r io.Reader, t *T) error {
	err := binary.Read(r, binary.LittleEndian, t)
	if err != nil {
		return err
	}

	return nil
}

// ReadStructBE reads arbitrary types and structs from io.Readers (big endian).
func ReadStructBE[T interface{}](r io.Reader, t *T) error {
	err := binary.Read(r, binary.BigEndian, t)
	if err != nil {
		return err
	}

	return nil
}

func IntAbs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// CurrentPos returns the return value of s.Seek(0, io.SeekCurrent).
// A seek with offset 0 should, in theory, never fail.
func CurrentPos(s io.ReadSeeker) int64 {
	pos, _ := s.Seek(0, io.SeekCurrent)
	return pos
}

// CopyPartOfFileToFile basically does exactly what it says on the tin.
// Useful for copying chunks from large files into new smaller files.
func CopyPartOfFileToFile(dst, src *os.File, srcOffset, srcLength int64) error {
	_, err := src.Seek(srcOffset, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.CopyN(dst, src, srcLength)
	if err != nil {
		return err
	}

	return nil
}

// HashFileSHA1 rewinds a file's pointer to the beginning, then returns an SHA1 hash of its contents.
func HashFileSHA1(f *os.File) (string, error) {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", h.Sum(nil)), nil
}

// IndexOfSlice returns the index of the needle slice in the haystack slice, or -1 if haystack does not contain needle.
// If haystack is smaller than needle, -1 is returned.
// If both slices' lengths are equal, 0 or -1 is returned, depending on whether the slices are equal.
func IndexOfSlice[T comparable](haystack []T, needle []T) int {
	haystackLen := len(haystack)
	needleLen := len(needle)

	if haystackLen < needleLen {
		return -1
	} else if haystackLen == needleLen {
		if slices.Equal(haystack, needle) {
			return 0
		}
		return -1
	}

	for i := 0; i <= haystackLen-needleLen; i++ {
		substack := haystack[i : i+needleLen]
		if slices.Equal(substack, needle) {
			return i
		}
	}

	return -1
}

// ReadNullTerminatedString scans an io.ReadSeeker for a maximum of 256 bytes and
// returns the first null-terminated string it finds, along with its offset in the stream.
// Returned strings will not contain any null bytes.
func ReadNullTerminatedString(f io.ReadSeeker) (int64, string, error) {
	// TODO: strings are 8-byte aligned. implementing that might save us a bunch of seeking
	var tempBytes []byte
	var offset int64 = -1
	for i := 0; i < 256; i++ {
		tempByte := make([]byte, 1)
		_, err := f.Read(tempByte)
		if err != nil {
			if errors.Is(err, io.EOF) && tempBytes != nil {
				return offset, string(tempBytes), nil
			}

			return -1, "", err
		}

		b := tempByte[0]
		if b == 0 {
			// string terminator
			if tempBytes == nil {
				// we haven't gotten any good ASCII bytes yet.
				// continue until we get at least one
				continue
			}

			return offset, string(tempBytes), nil
		} else if b < 0x20 || b >= 0x7F {
			// invalid ASCII char
			return -1, "", ErrInvalidASCIIChar
		} else {
			// we found an ASCII char!
			// add it to tempBytes and store the offset we found it at
			tempBytes = append(tempBytes, b)
			if offset < 0 {
				pos, _ := f.Seek(0, io.SeekCurrent)
				offset = pos - 1 // minus one to compensate for the read earlier
			}
		}
	}

	return -1, "", ErrMaxLengthExceeded
}
