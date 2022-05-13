package gbook

import (
	"bytes"
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

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
	Https       bool
}

// Info is the output book's information
type Info struct {
	Title            string
	SourceFolderPath string
	OutputFolderPath string
	menuContent      string
	// Whether generate a check file
	CompileCheckFile bool

	Template     fs.ReadFileFS
	TemplatePath string

	ServeInfo
}

const (
	indexPage = "README.md"
	menuPage  = "SUMMARY.md"
)

func New() *Info {
	return &Info{
		Title:            "GBook",
		Template:         defaultTemplate,
		TemplatePath:     "templates/default",
		CompileCheckFile: true,
		ServeInfo: ServeInfo{
			Port: "4000",
		},
	}
}

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

func (i *Info) InitOutputFolder() {
	i.generateOutputFolder()
	log.Printf("Init the output folder %s", i.OutputFolderPath)
	copy.FS(i.Template, i.TemplatePath, i.OutputFolderPath)

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
	filesOfInputFolder, err := sio.LoadFilesInFolderIgnoreHiddenFiles(i.SourceFolderPath)
	if err != nil {
		return
	}
	for index, f := range filesOfInputFolder {
		p, _ := filepath.Rel(i.SourceFolderPath, f)

		filesOfInputFolder[index] = p
	}
	i.generateMenu()
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
			i.generatePage(relativeFilePath)
		} else {
			copy.File(i.SourceFolderPath, i.OutputFolderPath, relativeFilePath)
		}
	}

	return
}

func (i *Info) generatePage(relativeFilePath string) {
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
	mdData = bytes.ReplaceAll(mdData, []byte("\r\n"), []byte("\n"))

	urlPrefix, _ := filepath.Rel("/"+filepath.Dir(relativeFilePath), "/")
	urlPrefix = strings.ReplaceAll(path.Join(urlPrefix, "/"), `\`, "/")

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
				Flags: html.CommonFlags | html.NoreferrerLinks | html.NoopenerLinks,
			},
		),
	)

	outFile, err := os.Create(path.Join(i.OutputFolderPath, outputRelativeFilePath))
	if err != nil {
		return
	}
	defer outFile.Close()
	if i.CompileCheckFile {
		ioutil.WriteFile(path.Join(i.OutputFolderPath, strings.Replace(relativeFilePath, ".md", ".compiled.xml", 1)), outContent, 0644)
	}
	t := template.Must(template.ParseFS(defaultTemplate, "templates/default/*.tmpl"))
	t.ExecuteTemplate(outFile, "index.tmpl", struct{ Title, Menu, MainContent, Prefix interface{} }{
		Title:       i.Title,
		Menu:        strings.ReplaceAll(i.menuContent, `href="/`, `href="`+urlPrefix+"/"),
		MainContent: string(outContent),
		Prefix:      urlPrefix,
	})
}

// generate the menu of nav-bar
func (i *Info) generateMenu() {
	log.Println("Generating the menu")
	defer log.Print("Finish generating the menu")

	mdData, err := ioutil.ReadFile(path.Join(i.SourceFolderPath, menuPage))
	if err != nil {
		i.menuContent = "<!-- " + err.Error() + "-->"
		return
	}

	parentNode := markdown.Parse(mdData, parser.New())
	ast.WalkFunc(parentNode, handleLinkTag)

	i.menuContent = string(
		markdown.Render(
			parentNode,
			html.NewRenderer(
				html.RendererOptions{
					Flags: html.CommonFlags,
				},
			),
		),
	)
	i.menuContent = strings.ReplaceAll(strings.ReplaceAll(i.menuContent, "<ul>", "<ol>"), "</ul>", "</ol>")
}

func handleLinkTag(node ast.Node, entering bool) ast.WalkStatus {
	if link, ok := node.(*ast.Link); ok {
		if match, err := regexp.Match(`\b([A-Za-z]+:|)//.+`, link.Destination); err != nil {
			log.Print(err)
		} else if match {
			for _, t := range link.AdditionalAttributes {
				if t == `target="_blank"` {
					return ast.GoToNext
				}
			}
			link.AdditionalAttributes = append(link.AdditionalAttributes, `target="_blank"`)
		} else {

			u := string(link.Destination)
			dir, file := filepath.Split(u)
			if file == indexPage {
				link.Destination = []byte(path.Join(dir, "index.html"))
			} else if file[len(file)-3:] == ".md" {
				link.Destination = []byte(u[0:len(u)-3] + ".html")
			}
		}
	}

	return ast.GoToNext
}
