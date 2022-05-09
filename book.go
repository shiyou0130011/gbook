package gbook

import (
	"embed"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	sio "github.com/shiyou0130011/io"
	"github.com/shiyou0130011/io/copy"
)

//go:embed templates/default
var defaultTemplate embed.FS

//ServeInfo is the serving data of the book
type ServeInfo struct {
	Port        string
	CertKeyPath string
}

// Info is the output book's information
type Info struct {
	Title            string
	SourceFolderPath string
	OutputFolderPath string

	ServeInfo
}

const (
	indexPage = "README.md"
	menuPage  = "SUMMARY.md"
)

// To generate the output folder.
// if the output folder is not exist, it will also create it.
func (i *Info) generateOutputFolder() (err error) {
	if i.OutputFolderPath == "" {
		i.OutputFolderPath, err = ioutil.TempDir(os.TempDir(), "gbook-"+path.Base(strings.ReplaceAll(i.SourceFolderPath, "\\", "/"))+"-*")
		if err == nil {
			return
		}
	}

	if _, fileExistErr := os.Stat(i.OutputFolderPath); os.IsNotExist(fileExistErr) {
		log.Print("Create folder ", i.OutputFolderPath)
		err = os.Mkdir(i.OutputFolderPath, os.ModeDir)
	}

	return
}

func (i *Info) InitOutputFolderWithFS(template fs.ReadFileFS, templatePath string) {
	i.generateOutputFolder()
	log.Printf("Init the output folder %s", i.OutputFolderPath)
	copy.FS(template, templatePath, i.OutputFolderPath)

	fs.WalkDir(os.DirFS(i.OutputFolderPath), ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".tmpl" {
			fileFullPath := filepath.Join(i.OutputFolderPath, path)
			log.Printf("Remove %s", fileFullPath)
			os.Remove(fileFullPath)
		}
		return nil
	})
}
func (i *Info) Compile() (err error) {
	filesOfInputFolder, err := sio.LoadFilesInFolder(i.SourceFolderPath)
	if err != nil {
		return
	}
	for index, f := range filesOfInputFolder {
		p, _ := filepath.Rel(i.SourceFolderPath, f)

		filesOfInputFolder[index] = p
	}
	menuContent := ""
	for _, relativeFilePath := range filesOfInputFolder {
		if relativeFilePath == menuPage {
			continue
		}

		outDir := path.Join(i.OutputFolderPath, filepath.Dir(relativeFilePath))
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			log.Print("Create folder ", outDir)
			os.Mkdir(outDir, os.ModeDir)
		}

		if path.Ext(relativeFilePath) == ".md" {
			i.generatePage(menuContent, relativeFilePath)
		} else {
			copy.File(i.SourceFolderPath, i.OutputFolderPath, relativeFilePath)
		}
	}

	return
}

func (i *Info) generatePage(menuContent string, relativeFilePath string) {
	dir, f := filepath.Split(relativeFilePath)
	var outputRelativeFilePath string
	if f == "README.md" {
		outputRelativeFilePath = path.Join(dir, "index.html")
	} else {
		outputRelativeFilePath = path.Join(dir, f[0:len(f)-2]+"html")
	}

	log.Printf("Generating %s (output: %s)", relativeFilePath, outputRelativeFilePath)

	mdData, err := ioutil.ReadFile(path.Join(i.SourceFolderPath, relativeFilePath))
	if err != nil {
		return
	}

	outNode := markdown.Parse(
		mdData,
		parser.NewWithExtensions(
			parser.CommonExtensions|parser.AutoHeadingIDs,
		),
	)
	//ast.WalkFunc(outNode, handleLinkTag)
	outContent := markdown.Render(
		outNode,
		html.NewRenderer(
			html.RendererOptions{
				//Flags: html.CommonFlags | html.TOC | html.CommonFlags,
			},
		),
	)

	outFile, err := os.Create(path.Join(i.OutputFolderPath, outputRelativeFilePath))
	if err != nil {
		return
	}
	defer outFile.Close()
	ioutil.WriteFile(path.Join(i.OutputFolderPath, relativeFilePath), outContent, 0644)

	t := template.Must(template.ParseFS(defaultTemplate, "templates/default/*.tmpl"))
	t.ExecuteTemplate(outFile, "index.tmpl", struct{ Title, Menu, MainContent interface{} }{
		Title:       i.Title,
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
