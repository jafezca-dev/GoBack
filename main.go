package main

import (
	"os"
)

func getFiles(path string, files *[]FileDiff) {
	content, _ := os.ReadDir(path)

	//if error != nil {
	//	panic(error)
	//}

	for _, file := range content {
		if file.IsDir() {
			getFiles(path+"\\"+file.Name(), files)
		} else {
			*files = append(*files, FileDiff{NewFile: file, FullDirPath: path})
		}
	}
}

func main() {
	var files []FileDiff

	const BasePath = "C:\\Users\\Javi\\Documents\\Isos"

	getFiles(BasePath, &files)

	//diffs := map[string]FileDiff{}

	//for _, file := range files {
	//	_, fileName := file.DirPaths(BasePath)
	//
	//	diffs[fileName] = file
	//}
}
