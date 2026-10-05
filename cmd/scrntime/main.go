package main

import (
	"github.com/Akhil373/kde-screentime/internal/cli"

	"github.com/alexflint/go-arg"
)

type Args struct {
	Graph     bool `arg:"-g,--graph" help:"render the graph"`
	NumOfDays int  `arg:"positional" default:"7" help:"number of days of data to fetch"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.Graph {
		cli.RunGraph(args.NumOfDays)
		return
	}
	cli.RunDaemon()
}
