package main

import (
	"flag"
	"gbook"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	indexPage = "README.md"
	menuPage  = "SUMMARY.md"
)

func main() {
	book := gbook.New()

	var (
		showHelp  bool
		serveHTTP bool
	)
	flag.StringVar(&book.SourceFolderPath, "f", ".", "The folder of the book")
	flag.StringVar(&book.Title, "title", "GBook", "The title of the book")
	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.BoolVar(&serveHTTP, "serve", false, "Serve the book")
	flag.StringVar(&book.OutputFolderPath, "out", "", "The output folder path. If it is blank, the program will create a new folder and set the parameter as the folder's path. ")
	flag.StringVar(&book.Port, "p", "4000", "When serving the book, the HTTP port for serving site")
	flag.StringVar(&book.CertKeyPath, "c", "", "When serving the book, the folder used for HTTPS cert key path")

	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		return
	}

	// check input folder
	_, err := os.Open(path.Join(book.SourceFolderPath, indexPage))
	if err != nil {
		log.Fatalf("Cannot find index page (%s)", path.Join(book.SourceFolderPath, indexPage))
	}

	// generate output folder
	book.InitOutputFolder()
	book.Compile()

	if serveHTTP {
		if book.CertKeyPath != "" {
			http.ListenAndServeTLS(
				":"+book.Port,
				filepath.Join(book.CertKeyPath, "cert.pem"),
				filepath.Join(book.CertKeyPath, "key.pem"),
				http.FileServer(http.Dir(book.OutputFolderPath)),
			)
		} else {
			http.ListenAndServe(":"+book.Port, http.FileServer(http.Dir(book.OutputFolderPath)))
		}

	}
}
