package bin

// not sure what this is for yet.
// this *might* help find the data table offset in other SH2 binary versions

const paddingLength = 256

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
