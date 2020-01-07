package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vyh/geoipupdate/pkg/geoipupdate"
	"github.com/vyh/geoipupdate/pkg/geoipupdate/database"
)

var (
	version                  = "unknown"
	defaultConfigFile        string
	defaultDatabaseDirectory string
)

func main() {
	log.SetFlags(0)

	if defaultConfigFile == "" {
		defaultConfigFile = geoipupdate.DefaultConfigFile
	}
	if defaultDatabaseDirectory == "" {
		defaultDatabaseDirectory = geoipupdate.DefaultDatabaseDirectory
	}

	args := getArgs()
	fatalLogger := func(message string, err error) {
		if args.StackTrace {
			log.Print(fmt.Sprintf("%s: %+v", message, err))
		} else {
			log.Print(fmt.Sprintf("%s: %s", message, err))
		}
		os.Exit(1)
	}

	config, err := geoipupdate.NewConfig(args.ConfigFile, defaultDatabaseDirectory, args.DatabaseDirectory, args.Verbose)
	if err != nil {
		fatalLogger("error loading configuration file", err)
	}
	if config.Verbose {
		log.Printf("Using config file %s", args.ConfigFile)
		log.Printf("Using database directory %s", config.DatabaseDirectory)
	}

	if err = run(config); err != nil {
		fatalLogger("error retrieving updates", err)
	}
}

func run(
	config *geoipupdate.Config,
) error {
	client := geoipupdate.NewClient(config)
	dbReader := database.NewHTTPDatabaseReader(client, config)

	for _, editionID := range config.EditionIDs {
		filename, err := geoipupdate.GetFilename(config, editionID, client)
		if err != nil {
			return errors.Wrap(err, "error retrieving filename")
		}
		filePath := filepath.Join(config.DatabaseDirectory, filename)
		dbWriter, err := database.NewLocalFileDatabaseWriter(filePath, config.LockFile, config.Verbose)
		if err != nil {
			return errors.Wrap(err, "error creating database writer")
		}
		if err := dbReader.Get(dbWriter, editionID); err != nil {
			return err
		}
	}
	return nil
}
