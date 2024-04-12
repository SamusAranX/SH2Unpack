# SH2Unpack

This tool is a Silent Hill 2 file extractor and was written to help a friend with a video about the game.

I hope it can be useful for other people interested in SH2 as well.

## Downloads

The latest release can be found here: https://github.com/SamusAranX/sh2unpack/releases/latest

## Usage

This is a commandline tool, which means you'll be running it in the terminal, cmd.exe, or your platform equivalent. 
Double-clicking or dropping files on it will do nothing.

**NOTE: This tool operates on files from the PS2 version of the game. It was not tested with other versions, which are beyond its scope.**
The only specific version of the game this tool was developed around is the NTSC Greatest Hits version, also known as v2.01.
Other versions haven't been tested yet.

Also, you'll have to copy all files from a game ISO to another folder. The example below assumes the files are in a folder called "SH2". 

This tool currently has one command `unpack` that takes one argument `-i <input file>` and one last argument that's the output directory.
That's where extracted files go.

Here's an example:

```
$ sh2unpack unpack -i ./SH2/SLUS_202.28 ./SH2Unpack/
```

The tool will output something like:

```
sh2unpack v1.0-dirty [10c3baa]
Input: <path>/SH2/SLUS_202.28
Output Folder: <path>/SH2Unpack/
Version detected: SLUS_202.28, NTSC v2.01 (Greatest Hits)
Extracted 3825 files.
```

## Supported game versions

This tool currently supports 10 distinct versions of the game.\
This includes the entire set of discs known to redump.org
(minus the Tentou Houei-you Movie-ban disc on account of it containing no extractable data) plus the E3 2001 demo.

| Binary Name | Region      | Version Info           | SHA1 Hash (ISO File)                       | SHA1 Hash (Binary)                         |
|-------------|-------------|------------------------|--------------------------------------------|--------------------------------------------|
| SLPM_123.45 | NTSC-U 🇺🇸 | E3 2001                | `02F2E34E018596A31C0A5CAB1B6BA981ABC2F008` | `50C664C525736619215654186446A5D6B211FB31` |
| SLUS_202.28 | NTSC-U 🇺🇸 | Jul 13, 2001 prototype | `BBEBD65FCD3E792C3A57DBADF3EE1DEB2846172E` | `888EFF71606FF4C1C610E30111B3CA5DA647EDCC` |
| SLUS_202.28 | NTSC-U 🇺🇸 | Silent Hill 2          | `F7FCB40D8C79A6AC622299069A4DC2900C74B200` | `ECFD22C67F7712480F52D0674B70964D2A82E648` |
| SLUS_202.28 | NTSC-U 🇺🇸 | Greatest Hits          | `2F4D89736D9240C6F8719E50A8D450A81AD638AE` | `3A27DEDDFA81CF30F46F0742C3523230CAC75D9A` |
| SLES_503.82 | PAL 🇪🇺    | Special 2 Disc Set     | `924409DE4DC4CABD4A978FAE7DE94159E57A1C8D` | `8BC367E1B9E7AA5CC5D5FA32048ED97F3FADE728` |
| SLES_511.56 | PAL 🇪🇺    | Director's Cut         | `3A2B03AEF487AE88BA5C51B064AAF8295398F684` | `2C5A7AFBA3A5B4507CCB828811C8ADD9E5D0E961` |
| SLKA_250.01 | NTSC-J 🇰🇷 | Silent Hill 2          | `5A215C62899F4DA374B9F6E0E56CA6CA6D0A06CB` | `279A1B4DBFD43FF7A5920A52D51B153C638D1D6B` |
| SLPM_610.09 | NTSC-J 🇯🇵 | Red Ribbon Demo        | `469DDB3E50EEFBF2C5BBC39E1FDF6FC039AD502B` | `B9CB2E895FC83CD4452DC9A818BF3CA26394ADBE` |
| SLPM_650.51 | NTSC-J 🇯🇵 | Silent Hill 2          | `6A9C80C3D965EE0E50A4FC131AE0D3D9F2384552` | `ED1DB66E92FEE366B375D5A1993F4609641BE6DA` |
| SLPM_650.98 | NTSC-J 🇯🇵 | Saigo no Uta           | `9BDF3E49F22366B0C27EC9F1EE31721A5106B1B4` | `EFA89AA35054A9A547F22673AB601CFB333587DE` |

If there are other versions of the game you think this tool should support, please file an issue.

Modded versions are not *and will not be* officially supported.
A way to skip the hash recognition step and manually provide offsets will be implemented at a later date.

## Building

Use the makefile to create builds.

Specifically, `make build` to create a build for your current platform or `make buildall` to create 
builds for Windows (x64), macOS (x64, arm64), and Linux (x64, arm64).

## Helpful Links

To work with PSS video files, download PSS_demux v1.05 from this website.\
The website's in Japanese, but the tool's in English and **very** easy to use:\
https://azuco.sakura.ne.jp/fao/fao.html