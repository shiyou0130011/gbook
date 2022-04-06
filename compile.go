package main

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed templates/default
var defaultTemplate embed.FS

func compileMarkdownFiles(sourceFolderPath string, outputFolderPath string, bookTitle string) {
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
			generatePage(sourceFolderPath, outputFolderPath, menuContent, relativeFilePath, bookTitle)
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

	parentNode := markdown.Parse(mdData, parser.New())
	ast.WalkFunc(parentNode, handleLinkTag)

	return string(
		markdown.Render(
			parentNode,
			html.NewRenderer(
				html.RendererOptions{
					Flags: html.CommonFlags,
				},
			),
		),
	)
}

// Generate page
func generatePage(sourceFolderPath string, outputFolderPath string, menuContent string, relativeFilePath string, bookTitle string) {
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

	outNode := markdown.Parse(
		mdData,
		parser.NewWithExtensions(
			parser.CommonExtensions|parser.AutoHeadingIDs,
		),
	)
	ast.WalkFunc(outNode, handleLinkTag)
	outContent := markdown.Render(
		outNode,
		html.NewRenderer(
			html.RendererOptions{
				//Flags: html.CommonFlags | html.TOC | html.CommonFlags,
			},
		),
	)

	outFile, err := os.Create(path.Join(outputFolderPath, outputRelativeFilePath))
	if err != nil {
		return
	}
	defer outFile.Close()
	ioutil.WriteFile(path.Join(outputFolderPath, relativeFilePath), outContent, 0644)

	t := template.Must(template.ParseFS(defaultTemplate, "templates/default/*.html"))
	t.ExecuteTemplate(outFile, "index.html", struct{ Title, Menu, MainContent interface{} }{
		Title:       bookTitle,
		Menu:        menuContent,
		MainContent: string(outContent),
	})

}

func handleLinkTag(node ast.Node, entering bool) ast.WalkStatus {
	if link, ok := node.(*ast.Link); ok {

		if match, err := regexp.Match(`\b([A-Za-z]+:|)//.+`, link.Destination); err != nil {
			log.Print(err)
		} else if match {
			if link.Attribute == nil {
				link.Attribute = &ast.Attribute{
					Attrs: map[string][]byte{},
				}
			}
			link.Attrs["target"] = []byte("_blank")
		} else {

			u := string(link.Destination)
			dir, file := filepath.Split(u)
			if file == indexPage {
				link.Destination = []byte(filepath.Join(dir, "index.html"))
			} else if file[len(file)-3:] == ".md" {
				link.Destination = []byte(u[0:len(u)-3] + ".html")
			}
		}
	}

	return ast.GoToNext
}
