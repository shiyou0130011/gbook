package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	indexPage = "README.md"
	menuPage  = "SUMMARY.md"
)

func main() {
	var (
		sourceFolderPath string
		showHelp         bool
		serveHTTP        bool
		port             string
		title            string
		outFolderName    string
	)
	flag.StringVar(&sourceFolderPath, "f", ".", "The folder of the book")
	flag.StringVar(&title, "title", "GBook", "The title of the book")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.BoolVar(&serveHTTP, "serve", false, "Serve the book")
	flag.StringVar(&outFolderName, "out", "", "The output folder path. If it is blank, the program will create a new folder and set the parameter as the folder's path. ")
	flag.StringVar(&port, "p", "4000", "When serving the book, the HTTP port for serving site")

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
	if outFolderName == "" {
		outFolderName, err = ioutil.TempDir(os.TempDir(), "gbook-"+path.Base(strings.ReplaceAll(sourceFolderPath, "\\", "/"))+"-*")
		if err != nil {
			log.Fatal("Cannot generate output folder")
		}
		log.Printf("Generate output folder: %s", outFolderName)
	}

	compileMarkdownFiles(sourceFolderPath, outFolderName, title)

	if serveHTTP {
		http.ListenAndServe(":4000", http.FileServer(http.Dir(outFolderName)))
	}
}
