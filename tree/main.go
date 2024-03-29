package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

// FolderLikeLeaf - спомощью этой структуры создается связанный список папок
type FolderLikeLeaf struct {
	Name        string
	Path        string
	ChildFolder []*FolderLikeLeaf
	ChildItem   []os.FileInfo
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
	path := flag.String("path", "..", "the root leaf for tree")
	fileShow := flag.Bool("f", false, "show files")
	flag.Parse()

	out := new(bytes.Buffer)
	err := dirTree(out, *path, *fileShow)
	if err != nil {
		panic(err)
	}
	fmt.Print(out.String())
}

func dirTree(out *bytes.Buffer, myPath string, fileShow bool) error {
	folder, err := getListFolder(myPath, nil, fileShow)
	if err != nil {
		return err
	}
	*out = *getByteListFolder(folder, "")
	return nil
}

func getListFolder(myPath string, parent *FolderLikeLeaf, fileShow bool) (*FolderLikeLeaf, error) {
	stat, err := os.Stat(myPath)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() != true {
		return nil, fmt.Errorf("This is a file. The program expects a folder")
	}

	folder := &FolderLikeLeaf{Name: stat.Name(), Path: myPath,
		ChildFolder: nil}
	folders, err := ioutil.ReadDir(myPath)
	if err != nil {
		return nil, err
	}
	for _, item := range folders {
		if item.IsDir() == true {
			childFolder, err := getListFolder(path.Join(myPath, item.Name()), folder, fileShow)
			if err != nil {
				return folder, err
			}
			folder.ChildFolder = append(folder.ChildFolder, childFolder)
			folder.ChildItem = append(folder.ChildItem, item)
		}
		if item.IsDir() != true {
			if fileShow {
				folder.ChildItem = append(folder.ChildItem, item)
			}
		}
	}

	return folder, err
}

//rudiment
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
}

//rudiment
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

//rudiment
func printStringListFolder(listFolder *FolderLikeLeaf, tab string) string {
	var outStr string
	listFolder.sortChildItem()
	n := len(listFolder.ChildItem)
	for i, item := range listFolder.ChildItem {
		if i == n-1 {
			outStr = outStr + tab + "└───" + item.Name() + "\n"
			if item.IsDir() == true {
				outStr = outStr + printStringListFolder(listFolder.getChildFolder(item.Name()), tab+"	")
			}
		} else {
			outStr = outStr + tab + "├───" + item.Name() + "\n"
			if item.IsDir() == true {
				outStr = outStr + printStringListFolder(listFolder.getChildFolder(item.Name()), tab+"│	")
			}
		}
	}
	return outStr
}

func getByteListFolder(listFolder *FolderLikeLeaf, tab string) *bytes.Buffer {
	var tempStr string
	outBuffer := new(bytes.Buffer)
	listFolder.sortChildItem()
	n := len(listFolder.ChildItem)
	for i, item := range listFolder.ChildItem {
		if i == n-1 {
			if item.IsDir() == true {
				outBuffer.WriteString(tab + "└───" + item.Name() + "\n")
				outBuffer.Write(getByteListFolder(listFolder.getChildFolder(item.Name()), tab+"	").Bytes())
			} else {
				if item.Size() == 0 {
					tempStr = " (empty)\n"
				} else {
					tempStr = fmt.Sprintf(" (%db)\n", item.Size())
				}
				outBuffer.WriteString(tab + "└───" + item.Name() + tempStr)
			}
		} else {
			if item.IsDir() == true {
				outBuffer.WriteString(tab + "├───" + item.Name() + "\n")
				outBuffer.Write(getByteListFolder(listFolder.getChildFolder(item.Name()), tab+"│	").Bytes())
			} else {
				if item.Size() == 0 {
					tempStr = " (empty)\n"
				} else {
					tempStr = fmt.Sprintf(" (%db)\n", item.Size())
				}
				outBuffer.WriteString(tab + "├───" + item.Name() + tempStr)
			}
		}
	}
	return outBuffer
}
