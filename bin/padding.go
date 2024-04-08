package bin

const paddingLength = 256

// PS2Padding constructs and returns the 1024-byte block of padding found in Silent Hill 2's main binary.
func PS2Padding() []byte {
	// 0+1+2+4+8+16+32+64+128
	padding := make([]byte, paddingLength*4)

	for i := 0; i < 1; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 0
		}
	}
	for i := 1; i < 1+1; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 1
		}
	}
	for i := 2; i < 2+2; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 2
		}
	}
	for i := 4; i < 4+4; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 3
		}
	}
	for i := 8; i < 8+8; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 4
		}
	}
	for i := 16; i < 16+16; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 5
		}
	}
	for i := 32; i < 32+32; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 6
		}
	}
	for i := 64; i < 64+64; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 7
		}
	}
	for i := 128; i < 128+128; i++ {
		for j := 0; j < 4; j++ {
			padding[i+(j*paddingLength)] = 8
		}
	}

	return padding
}
