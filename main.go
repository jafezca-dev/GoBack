package main

import (
	"fmt"
	"os"
)

func getParameters(commandParams []string) ProgParams {
	progParams := ProgParams{}

	for index, param := range commandParams {
		switch param {
		case "-b":
			progParams.Bucket = commandParams[index+1]
		case "-ep":
			progParams.Endpoint = commandParams[index+1]
		case "-ak":
			progParams.AccessKey = commandParams[index+1]
		case "-sk":
			progParams.SecretKey = commandParams[index+1]
		}
	}

	return progParams
}

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
	commandParams := os.Args[1:]
	progParams := getParameters(commandParams)
	fmt.Println(progParams)

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
