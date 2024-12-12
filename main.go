package main

import (
	"GoBack/clients"
	"GoBack/types"
	"bytes"
	"os"
	"time"
)

func getParameters(commandParams []string) types.ProgParams {
	progParams := types.ProgParams{}

	for index, param := range commandParams {
		switch param {
		case "full", "incremental":
			progParams.BackupType = param
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

func addOldFiles(progParams types.ProgParams, storageClient clients.StorageClient, files *[]types.FileDiff) {
	oldDataList, err := storageClient.GetLastChanges()
	if err == nil {
		for index, file := range *files {
			_, virtualFilePath := file.DirPaths(progParams.BasePath)
			fileOldData, exists := oldDataList[virtualFilePath]
			if exists {
				(*files)[index].OldFile = &fileOldData
			}
		}
	}
}

func main() {
	commandParams := os.Args[1:]
	progParams := getParameters(commandParams)

	storageClient := getStorageClient(progParams)
	storageClient.CheckBucketConnection()

	var files []types.FileDiff

	getFiles(progParams.BasePath, &files)

	if progParams.BackupType == "incremental" {
		addOldFiles(progParams, storageClient, &files)
	}

	var csvBuffer bytes.Buffer

	filesUploadeds := 0

	for _, file := range files {
		if storageClient.UploadFile(file) {
			filesUploadeds++
		}

		csvLine := file.GetCsvReg(progParams)
		_, err := csvBuffer.WriteString(csvLine + "\n")
		if err != nil {
			panic("Error CSV")
		}
	}

	if filesUploadeds != 0 {
		storageClient.UploadCsv(csvBuffer)
	}
}
