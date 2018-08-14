package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

// FolderLikeLeaf - спомощью этой структуры создается связанный список папок
type FolderLikeLeaf struct {
	Name         string
	Path         string
	ChildFolder  []*FolderLikeLeaf
	ChildFile    []string
	ChildItem    []os.FileInfo
	ParentFolder *FolderLikeLeaf
	RootFolder   *FolderLikeLeaf
}

func (folder FolderLikeLeaf) getChildFolder(name string) *FolderLikeLeaf {
	for _, item := range folder.ChildFolder {
		if name == item.Name {
			return item
		}
	}
	return nil
}

func (folder FolderLikeLeaf) sortChildItem() {
	sort.Slice(folder.ChildItem, func(i, j int) bool { return folder.ChildItem[i].Name() < folder.ChildItem[j].Name() })
}

func main() {
	/*	folder, err := getListFolder("..", nil)
		if err != nil {
			fmt.Println(err)
		}
		//	fmt.Println(folder)
		//printfListFolder(folder, "")
		fmt.Print(printStringListFolder(folder, ""))*/
	fmt.Print(dirTree(".."))
}

func dirTree(myPath string) string {
	folder, err := getListFolder(myPath, nil)
	if err != nil {
		fmt.Println(err)
	}
	return printStringListFolder(folder, "")
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
	folder.ChildItem = folders
	for _, item := range folders {
		if item.IsDir() == true {
			childFolder, err := getListFolder(path.Join(myPath, item.Name()), folder)
			if err != nil {
				return folder, err
			}
			folder.ChildFolder = append(folder.ChildFolder, childFolder)
		}
		if item.IsDir() != true {
			folder.ChildFile = append(folder.ChildFile, item.Name())
		}
	}

	return folder, err
}

func printListFolder(listFolder *FolderLikeLeaf, tab string) {
	fmt.Println(listFolder.Name)
	listFolder.sortChildItem()
	for _, item := range listFolder.ChildItem {
		fmt.Print(tab)
		if item.IsDir() != true {
			fmt.Println(item.Name())
		}
		if item.IsDir() == true {
			printListFolder(listFolder.getChildFolder(item.Name()), tab+"\t")
		}
	}
	/*	for _, item := range listFolder.ChildFolder {
		fmt.Print(tab)
		printListFolder(item, tab+"\t")
	}*/
}

func printfListFolder(listFolder *FolderLikeLeaf, tab string) {
	listFolder.sortChildItem()
	n := len(listFolder.ChildItem)
	for i, item := range listFolder.ChildItem {
		if i == n-1 {
			fmt.Print(tab)
			fmt.Print("└───")
			fmt.Println(item.Name())
			if item.IsDir() == true {
				printfListFolder(listFolder.getChildFolder(item.Name()), tab+"	")
			}
		} else {
			fmt.Print(tab)
			fmt.Print("├───")
			fmt.Println(item.Name())
			if item.IsDir() == true {
				printfListFolder(listFolder.getChildFolder(item.Name()), tab+"│	")
			}
		}
	}
}

func printStringListFolder(listFolder *FolderLikeLeaf, tab string) string {
	var outStr string
	listFolder.sortChildItem()
	n := len(listFolder.ChildItem)
	for i, item := range listFolder.ChildItem {
		if i == n-1 {
			outStr = outStr + tab + "└───" + item.Name() + "\n"
			//	fmt.Print(tab)
			//	fmt.Print("└───")
			//	fmt.Println(item.Name())
			if item.IsDir() == true {
				outStr = outStr + printStringListFolder(listFolder.getChildFolder(item.Name()), tab+"	")
			}
		} else {
			//	fmt.Print(tab)
			//	fmt.Print("├───")
			//	fmt.Println(item.Name())
			outStr = outStr + tab + "├───" + item.Name() + "\n"
			if item.IsDir() == true {
				outStr = outStr + printStringListFolder(listFolder.getChildFolder(item.Name()), tab+"│	")
			}
		}
	}
	return outStr
}
