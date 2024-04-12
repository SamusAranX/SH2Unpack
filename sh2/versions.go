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
		// NTSC-U
		"ECFD22C67F7712480F52D0674B70964D2A82E648": {
			DataOffset:  0x2BB180,
			MagicOffset: 0xFF900,
			FileName:    "SLUS_202.28",
			Description: "Silent Hill 2 (NTSC-U)",
		},
		"3A27DEDDFA81CF30F46F0742C3523230CAC75D9A": {
			DataOffset:  0x2CCF00,
			MagicOffset: 0xFF800,
			FileName:    "SLUS_202.28",
			Description: "Greatest Hits (NTSC-U)",
		},

		// NTSC-J
		"ED1DB66E92FEE366B375D5A1993F4609641BE6DA": {
			DataOffset:  0x2BB900,
			MagicOffset: 0xFF900,
			FileName:    "SLPM_650.51",
			Description: "Silent Hill 2 (NTSC-J, Japan)",
		},
		"279A1B4DBFD43FF7A5920A52D51B153C638D1D6B": {
			DataOffset:  0x2CD080,
			MagicOffset: 0xFF800,
			FileName:    "SLKA_250.01",
			Description: "Silent Hill 2 (NTSC-J, South Korea)",
		},
		"EFA89AA35054A9A547F22673AB601CFB333587DE": {
			DataOffset:  0x2CCB80,
			MagicOffset: 0xFF800,
			FileName:    "SLPM_650.98",
			Description: "Saigo no Uta (NTSC-J)",
		},

		// PAL
		"8BC367E1B9E7AA5CC5D5FA32048ED97F3FADE728": {
			DataOffset:  0x2BD400,
			MagicOffset: 0xFF800,
			FileName:    "SLES_503.82",
			Description: "Special 2 Disc Set (PAL)",
		},
		"2C5A7AFBA3A5B4507CCB828811C8ADD9E5D0E961": {
			DataOffset:  0x2CD980,
			MagicOffset: 0xFF800,
			FileName:    "SLES_511.56",
			Description: "Director's Cut (PAL)",
		},

		// Demos/Prototypes
		"50C664C525736619215654186446A5D6B211FB31": {
			DataOffset:  0x45C200,
			MagicOffset: 0xFFF80,
			FileName:    "SLPM_123.45",
			Description: "E3 2001 (NTSC-U)",
		},
		"888EFF71606FF4C1C610E30111B3CA5DA647EDCC": {
			DataOffset:  0x29CD00,
			MagicOffset: 0xFF900,
			FileName:    "SLUS_202.28",
			Description: "Jul 13, 2001 prototype (NTSC-U)",
		},
		"B9CB2E895FC83CD4452DC9A818BF3CA26394ADBE": {
			DataOffset:  0x2B3120,
			MagicOffset: 0xFF900,
			FileName:    "SLPM_610.09",
			Description: "Red Ribbon Demo (NTSC-J)",
		},
	}
)
