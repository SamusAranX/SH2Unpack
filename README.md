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

## Building

Use the makefile to create builds.

Specifically, `make build` to create a build for your current platform or `make buildall` to create builds for Windows (x64), macOS (x64, arm64), and Linux (x64, arm64).

## Helpful Links

To work with PSS video files, download PSS_demux v1.05 from this website.\
The website's in Japanese, but the tool's in English and **very** easy to use:\
https://azuco.sakura.ne.jp/fao/fao.html