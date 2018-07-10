package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// FolderLikeLeaf - спомощью этой структуры создается связанный список папок
type FolderLikeLeaf struct {
	Name         string
	Path         string
	ChildFolder  []*FolderLikeLeaf
	ChildFile    []string
	ParentFolder *FolderLikeLeaf
	RootFolder   *FolderLikeLeaf
}

func main() {
	fmt.Println("Hi World")
	folder, err := getListFolder(".", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(folder)
}

func dirTree(myPath string) string {
	return myPath
}

func getListFolder(myPath string, parent *FolderLikeLeaf) (*FolderLikeLeaf, error) {
	stat, err := os.Stat(myPath)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() != true {
		return nil, fmt.Errorf("This is a file. The program expects a folder")
	}

	//	var folder *FolderLikeLeaf
	folder := &FolderLikeLeaf{Name: stat.Name(), Path: myPath, ParentFolder: parent,
		ChildFolder: nil, ChildFile: nil}
	//	folder.Name = stat.Name()
	//	folder.Path = myPath
	//	folder.ParentFolder = parent
	folders, err := ioutil.ReadDir(myPath)
	if err != nil {
		return nil, err
	}
	for _, item := range folders {
		childFolder, err := getListFolder(path.Join(myPath, item.Name()), folder)
		if err != nil {
			return folder, err
		}
		folder.ChildFolder = append(folder.ChildFolder, childFolder)
	}

	return folder, err
}
