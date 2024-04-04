package bin

import (
	"errors"
	"io"
)

const (
	maxStringLength = 128
)

var (
	ErrInvalidASCIIChar  = errors.New("invalid ASCII char")
	ErrMaxLengthExceeded = errors.New("max search length exceeded")
	ErrWhaHappun         = errors.New("wha happun")
)

func readNullTerminatedString(f io.Reader) (string, error) {
	var tempBytes []byte
	for i := 0; true; i++ {
		if i == maxStringLength {
			return "", ErrMaxLengthExceeded
		}

		tempByte := make([]byte, 1)
		_, err := f.Read(tempByte)
		if err != nil {
			return "", err
		}

		b := tempByte[0]
		if b == 0 {
			// string terminator
			if tempBytes == nil {
				// we haven't gotten any good ASCII bytes yet.
				// continue until we get at least one
				continue
			}

			return string(tempBytes), nil
		} else if b < 0x20 || b >= 0x7F {
			// invalid ASCII char
			return "", ErrInvalidASCIIChar
		} else {
			// add to tempBytes
			tempBytes = append(tempBytes, b)
		}
	}

	return "", ErrWhaHappun
}
