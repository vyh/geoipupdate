package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Args are command line arguments.
type Args struct {
	ConfigFile        string `short:"f" description:"Configuration file (required)"`
	DatabaseDirectory string `short:"d" description:"Store databases in this directory (optional)"`
	StackTrace        bool   `long:"stack-trace" description:"Show a stack trace along with any error message"`
	Verbose           bool   `short:"v" description:"Use verbose output"`
	Version           bool   `short:"V" description:"Display the version and exit"`
}

func getArgs() *Args {
	args := &Args{}
	_, err := flags.Parse(args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if args.Version {
		log.Printf("geoipupdate %s", Version)
		os.Exit(0)
	}

	if *configFile == "" {
		log.Printf("You must provide a configuration file.")
		printUsage()
	}

	return args
}
