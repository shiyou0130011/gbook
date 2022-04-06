package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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

func copyFile(sourceFolderPath string, outputFolderPath string, relativeFilePath string) {
	log.Printf(`Copy file "%s" from "%s" to "%s"`, relativeFilePath, sourceFolderPath, outputFolderPath)
	data, err := ioutil.ReadFile(path.Join(sourceFolderPath, relativeFilePath))
	if err != nil {
		return
	}

	ioutil.WriteFile(path.Join(outputFolderPath, relativeFilePath), data, 0644)
}

//Copy the source folder to output folder
func copyDir(sourceFolder string, outputFolder string) {
	fileInfos, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		log.Printf("Cannot read folder %s", sourceFolder)
	}
	for _, fi := range fileInfos {
		sourcePath := filepath.Join(sourceFolder, fi.Name())
		outPath := filepath.Join(outputFolder, fi.Name())
		if fi.IsDir() {
			err = os.Mkdir(outPath, 0755)
			if err != nil {
				log.Printf("Cannot create folder %s", outPath)
			} else {
				copyDir(sourcePath, outPath)
			}

		} else {
			copyFile(sourcePath, outputFolder, fi.Name())
		}
	}
}
