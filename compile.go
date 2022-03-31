package main

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed templates/default
var defaultTemplate embed.FS

// load all files in folderPath
func loadFilesInFolder(folderPath string) (fileList []string) {
	files, err := ioutil.ReadDir(folderPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, loadFilesInFolder(path.Join(folderPath, file.Name()))...)
			continue
		}
		fileList = append(fileList, path.Join(folderPath, file.Name()))
	}
	return
}

func compileMarkdownFiles(sourceFolderPath string, outputFolderPath string) {
	filesOfInputFolder := loadFilesInFolder(sourceFolderPath)
	for index, f := range filesOfInputFolder {
		p, _ := filepath.Rel(sourceFolderPath, f)

		filesOfInputFolder[index] = p
	}

	menuContent := generateMenu(sourceFolderPath)

	for _, relativeFilePath := range filesOfInputFolder {
		if relativeFilePath == menuPage {
			continue
		}

		outDir := path.Join(outputFolderPath, filepath.Dir(relativeFilePath))
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			log.Print("Create folder ", outDir)
			os.Mkdir(outDir, os.ModeDir)
		}

		if path.Ext(relativeFilePath) == ".md" {
			generatePage(sourceFolderPath, outputFolderPath, menuContent, relativeFilePath)
		} else {
			copyFile(sourceFolderPath, outputFolderPath, relativeFilePath)
		}

	}
}

// generate the menu of nav-bar
func generateMenu(sourceFolderPath string) string {
	log.Println("Generating the menu")
	defer log.Print("Finish generating the menu")

	mdData, err := ioutil.ReadFile(path.Join(sourceFolderPath, menuPage))
	if err != nil {
		return "<!-- " + err.Error() + "-->"
	}
	return string(markdown.ToHTML(mdData, nil, nil))
}

func copyFile(sourceFolderPath string, outputFolderPath string, relativeFilePath string) {
	log.Printf(`Copy file "%s" from "%s" to "%s"`, relativeFilePath, sourceFolderPath, outputFolderPath)
	data, err := ioutil.ReadFile(path.Join(sourceFolderPath, relativeFilePath))
	if err != nil {
		return
	}

	ioutil.WriteFile(path.Join(outputFolderPath, relativeFilePath), data, 0644)
}

// Generate page
func generatePage(sourceFolderPath string, outputFolderPath string, menuContent string, relativeFilePath string) {
	dir, f := filepath.Split(relativeFilePath)
	var outputRelativeFilePath string
	if f == "README.md" {
		outputRelativeFilePath = path.Join(dir, "index.html")
	} else {
		outputRelativeFilePath = path.Join(dir, f[0:len(f)-2]+"html")
	}

	log.Printf("Generating %s (output: %s)", relativeFilePath, outputRelativeFilePath)

	mdData, err := ioutil.ReadFile(path.Join(sourceFolderPath, relativeFilePath))
	if err != nil {
		return
	}

	outContent := markdown.ToHTML(mdData,
		parser.NewWithExtensions(parser.CommonExtensions),
		html.NewRenderer(
			html.RendererOptions{
				Flags: html.CommonFlags | html.TOC | html.CommonFlags,
			},
		),
	)

	outFile, err := os.Create(path.Join(outputFolderPath, outputRelativeFilePath))
	if err != nil {
		return
	}
	defer outFile.Close()

	t := template.Must(template.ParseFS(defaultTemplate, "templates/default/*.html"))
	t.ExecuteTemplate(outFile, "index.html", struct{ Title, Menu, MainContent interface{} }{
		Title:       "G Book",
		Menu:        menuContent,
		MainContent: string(outContent),
	})

	// ioutil.WriteFile(, outContent, 0644)
}
