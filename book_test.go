package gbook

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_Info_Compile(t *testing.T) {
	inputFolder, _ := ioutil.TempDir(os.TempDir(), "input*")
	outputFolder, _ := ioutil.TempDir(os.TempDir(), "output*")
	defer os.RemoveAll(inputFolder)
	defer os.RemoveAll(outputFolder)

	content := `# Hello World

Lorem ipsum dolor sit amet, consectetur adipiscing elit`
	ioutil.WriteFile(inputFolder+"/README.md", []byte(content), 0644)

	book := New()
	book.SourceFolderPath = inputFolder
	book.OutputFolderPath = outputFolder
	book.CompileCheckFile = true

	book.InitOutputFolder()
	book.Compile()

	data, err := ioutil.ReadFile(outputFolder + "/README.compiled.xml")
	if err != nil {
		t.Fatal(err)

	}
	t.Logf("%s", data)

}

func Test_Info_CompileTable(t *testing.T) {
	inputFolder, _ := ioutil.TempDir(os.TempDir(), "input*")
	outputFolder, _ := ioutil.TempDir(os.TempDir(), "output*")
	defer os.RemoveAll(inputFolder)
	defer os.RemoveAll(outputFolder)

	content := `# Hello World

|Name    | Age |
|--------|------|
|Bob     | 27 |
|Alice   | 23|`
	ioutil.WriteFile(inputFolder+"/README.md", []byte(content), 0644)

	book := New()
	book.SourceFolderPath = inputFolder
	book.OutputFolderPath = outputFolder
	book.CompileCheckFile = true

	book.InitOutputFolder()
	book.Compile()

	data, err := ioutil.ReadFile(outputFolder + "/README.compiled.xml")
	if err != nil {
		t.Fatal(err)

	}
	t.Logf("%s", data)

}

func Test_Info_CompileList(t *testing.T) {
	inputFolder, _ := ioutil.TempDir(os.TempDir(), "input*")
	outputFolder, _ := ioutil.TempDir(os.TempDir(), "output*")
	defer os.RemoveAll(inputFolder)
	defer os.RemoveAll(outputFolder)

	content := `# Hello World

* Bob
* Alice

1. One 
2. Two`
	ioutil.WriteFile(inputFolder+"/README.md", []byte(content), 0644)

	book := New()
	book.SourceFolderPath = inputFolder
	book.OutputFolderPath = outputFolder
	book.CompileCheckFile = true

	book.InitOutputFolder()
	book.Compile()

	data, err := ioutil.ReadFile(outputFolder + "/README.compiled.xml")
	if err != nil {
		t.Fatal(err)

	}
	t.Logf("%s", data)

}
