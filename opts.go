package main

import "github.com/jessevdk/go-flags"

type DefaultOptions struct {
	Debug bool `long:"debug" description:"Debug mode"`
}

type UnpackOptions struct {
	DefaultOptions

	InFile flags.Filename `long:"infile" short:"i" required:"true" description:"The input file (.pdv or video file)"`

	Pos struct {
		OutDir flags.Filename `positional-arg-name:"outdir" description:"The output directory"`
	} `positional-args:"yes" required:"yes"`
}
