package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/drg"
)

type VersionFlag string

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Printf("DGR has version %s built from %s on %s\n", version, commit, date)
	app.Exit(0)

	return nil
}

var CLI struct {
	File string `arg:"" type:"existingfile" help:"Save file to manipulate"`
	Meta bool   `name:"meta" help:"Print metadata" optional:""`

	Version VersionFlag `name:"version" help:"Print version information and quit"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("DGR"),
		kong.UsageOnError(),
		kong.Description("Manipulate Deep Rock Galactic save files"),
	)

	ctx.FatalIfErrorf(run(ctx))
}

func run(ctx *kong.Context) error {
	f, err := os.Open(CLI.File)
	if err != nil {
		return err
	}
	defer f.Close()

	if CLI.Meta {
		meta, err := drg.DecodeMetadata(f)
		if err != nil {
			return err
		}

		jsonPrint(meta)
		return nil
	}

	save, err := drg.Decode(f)
	if err != nil {
		return err
	}
	jsonPrint(save)

	return nil
}

func jsonPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
