package main

import (
	"GoBack/clients"
	"GoBack/types"
	"os"
	"time"
)

func getParameters(commandParams []string) types.ProgParams {
	progParams := types.ProgParams{}

	for index, param := range commandParams {
		switch param {
		case "-p":
			progParams.BasePath = commandParams[index+1]
		case "-sc":
			progParams.StorageClient = commandParams[index+1]
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

	progParams.BackupDate = time.Now().Format("2006_01_02_15_04_05")

	return progParams
}

func getStorageClient(progParams types.ProgParams) clients.StorageClient {
	switch progParams.StorageClient {
	case "minio":
		return clients.NewMinioClient(progParams)
	}

	panic("No client selected")
}

func getFiles(path string, files *[]types.FileDiff) {
	content, _ := os.ReadDir(path)

	for _, file := range content {
		if file.IsDir() {
			getFiles(path+"\\"+file.Name(), files)
		} else {
			*files = append(*files, types.FileDiff{NewFile: file, FullDirPath: path})
		}
	}
}

func main() {

	//fmt.Println(backupDate)
	commandParams := os.Args[1:]
	progParams := getParameters(commandParams)

	storageClient := getStorageClient(progParams)
	storageClient.CheckBucketConnection()

	var files []types.FileDiff

	//const BasePath = "C:\\Users\\Javi\\Documents\\Isos"

	getFiles(progParams.BasePath, &files)

	//diffs := map[string]FileDiff{}

	for _, file := range files {
		storageClient.UploadFile(progParams, file)
	}

	//fmt.Println(diffs)
}
