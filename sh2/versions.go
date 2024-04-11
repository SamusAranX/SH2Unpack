package sh2

type gameVersion struct {
	DataOffset  uint32
	MagicOffset uint32
	FileName    string
	Description string
}

var (
	// map of game binary SHA1 hash -> game version
	VersionMap = map[string]gameVersion{
		// NTSC
		"3A27DEDDFA81CF30F46F0742C3523230CAC75D9A": {
			DataOffset:  0x2CCF00,
			MagicOffset: 0xFF800,
			FileName:    "SLUS_202.28",
			Description: "NTSC (Greatest Hits)",
		},

		// PAL
		"8BC367E1B9E7AA5CC5D5FA32048ED97F3FADE728": {
			DataOffset:  0x2BD400,
			MagicOffset: 0xFF800,
			FileName:    "SLES_503.82",
			Description: "PAL (Special 2 Disc Set)",
		},
		"2C5A7AFBA3A5B4507CCB828811C8ADD9E5D0E961": {
			DataOffset:  0x2CD980,
			MagicOffset: 0xFF800,
			FileName:    "SLES_511.56",
			Description: "PAL (Director's Cut)",
		},

		// Demos/Prototypes
		"50C664C525736619215654186446A5D6B211FB31": {
			DataOffset:  0x45C200,
			MagicOffset: 0xFFF80,
			FileName:    "SLPM_123.45",
			Description: "NTSC (E3 2001)",
		},
		"888EFF71606FF4C1C610E30111B3CA5DA647EDCC": {
			DataOffset:  0x29CD00,
			MagicOffset: 0xFF900,
			FileName:    "SLUS_202.28",
			Description: "NTSC (2001-07-13 prototype)",
		},
		"B9CB2E895FC83CD4452DC9A818BF3CA26394ADBE": {
			DataOffset:  0x2B3120,
			MagicOffset: 0xFF900,
			FileName:    "SLPM_610.09",
			Description: "PAL (Trial Version)",
		},
	}
)
