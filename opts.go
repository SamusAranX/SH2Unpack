package main

import "github.com/jessevdk/go-flags"

type DefaultOptions struct {
	Debug  bool `long:"debug" description:"Debug mode"`
	DryRun bool `long:"dry-run" description:"Skip file extraction"`
}

type UnpackOptions struct {
	DefaultOptions

	InFile flags.Filename `long:"infile" short:"i" required:"true" description:"The game's binary file (usually named something like SLUS_202.28)"`

	Pos struct {
		OutDir flags.Filename `positional-arg-name:"outdir" description:"The output directory"`
	} `positional-args:"yes" required:"yes"`
}
