package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	indexPage = "README.md"
	menuPage  = "SUMMARY.md"
)

func main() {
	sourceFolderPath := ""
	showHelp := false
	flag.StringVar(&sourceFolderPath, "f", ".", "The folder of the book")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		return
	}

	// check input folder
	_, err := os.Open(path.Join(sourceFolderPath, indexPage))
	if err != nil {
		log.Fatalf("Cannot find index page (%s)", path.Join(sourceFolderPath, indexPage))
	}

	_, err = os.Open(path.Join(sourceFolderPath, menuPage))
	if err != nil {
		log.Fatalf("Cannot find menu page (%s)", path.Join(sourceFolderPath, menuPage))
	}

	// generate output folder

	outFolderName, err := ioutil.TempDir(os.TempDir(), "gbook-"+path.Base(strings.ReplaceAll(sourceFolderPath, "\\", "/"))+"-*")
	if err != nil {
		log.Fatal("Cannot generate output folder")
	}
	log.Printf("Generate output folder: %s", outFolderName)

	outFolder, _ := os.Open(outFolderName)
	defer outFolder.Close()

	log.Printf("%#v", loadFilesInFolder(sourceFolderPath))
}
