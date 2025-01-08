package main

import (
	"GoBack/clients"
	"GoBack/types"
	"bytes"
	"os"
	"time"
)

func getParameters(commandParams []string) types.ProgParams {
	progParams := types.ProgParams{
		IgnoreFolders: map[string]bool{},
		IgnoreFiles:   map[string]bool{},
	}

	for index, param := range commandParams {
		switch param {
		case "full", "incremental":
			progParams.BackupType = param
		case "--path":
			progParams.BasePath = commandParams[index+1]
		case "--client":
			progParams.StorageClient = commandParams[index+1]
		case "--bucket":
			progParams.Bucket = commandParams[index+1]
		case "--endpoint":
			progParams.Endpoint = commandParams[index+1]
		case "--accesskey":
			progParams.AccessKey = commandParams[index+1]
		case "--secretkey":
			progParams.SecretKey = commandParams[index+1]
		case "--ignorefolder":
			progParams.IgnoreFolders[commandParams[index+1]] = true
		case "--ignorefile":
			progParams.IgnoreFiles[commandParams[index+1]] = true
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

func getFiles(path string, files *[]types.FileDiff, progParams types.ProgParams) {
	content, _ := os.ReadDir(path)

	for _, file := range content {
		if file.IsDir() && !progParams.IgnoreFolders[file.Name()] {
			getFiles(path+"\\"+file.Name(), files, progParams)
		} else if !progParams.IgnoreFiles[file.Name()] {
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

	getFiles(progParams.BasePath, &files, progParams)

	if progParams.BackupType == "incremental" {
		addOldFiles(progParams, storageClient, &files)
	}

	var csvBuffer bytes.Buffer

	//filesUploads := 0
	success := storageClient.MultiThreadUpload(files)

	for _, file := range files {
		csvLine := file.GetCsvReg(progParams)
		_, err := csvBuffer.WriteString(csvLine + "\n")
		if err != nil {
			panic("Error CSV")
		}
	}

	if success {
		storageClient.UploadCsv(csvBuffer)
	}
}
