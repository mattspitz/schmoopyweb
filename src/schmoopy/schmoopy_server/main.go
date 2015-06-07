package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	schmoopy ".."
)

const (
	required = "<REQUIRED>"
)

func main() {
	var dbFilename string
	var initDb bool
	var addr string
	var templateDir string
	var staticDir string

	flag.StringVar(&dbFilename, "db", required, "Filename for sqlite3 database.")
	flag.BoolVar(&initDb, "init", false, "If true, initializes the sqlite3 database")
	flag.StringVar(&addr, "addr", ":8080", "The TCP network address on which to run.")
	flag.StringVar(&templateDir, "templateDir", required, "The location of HTML templates for rendering schmoopies.")
	flag.StringVar(&staticDir, "staticDir", "", "The location of static files to be served at /static/ (optional)")

	flag.Parse()

	if dbFilename == required {
		log.Fatal(errors.New("-db is required!"))
	}

	if initDb {
		if err := schmoopy.InitializeDb(dbFilename); err != nil {
			log.Fatal(err)
		}
		return
	} else if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		log.Fatal(fmt.Errorf("File doesn't exist and we're not meant to initialize the database!: %v", dbFilename))
	}

	if templateDir == required {
		log.Fatal(errors.New("-templateDir is required!"))
	}

	s, err := schmoopy.NewSchmoopyServer(dbFilename, addr, templateDir, staticDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing server at", addr)
	log.Fatal(s.ListenAndServe())
}
