package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

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

	for _, filePath := range filesOfInputFolder {
		if filePath == menuPage {
			continue
		}

		outDir := path.Join(outputFolderPath, filepath.Dir(filePath))
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			log.Print("Create folder ", outDir)
			os.Mkdir(outDir, os.ModeDir)
		}

		// if path.Ext(filePath) == ".md" {
		// 	log.Print(filePath)
		// }

	}
}

// generate the menu of nav-bar
func generateMenu(sourceFolderPath string) ([]byte, error) {
	log.Println("Generating the menu")
	defer log.Print("Finish generating the menu")

	mdData, err := ioutil.ReadFile(path.Join(sourceFolderPath, menuPage))
	if err != nil {
		return nil, err
	}
	return markdown.ToHTML(mdData, nil, nil), nil
}
