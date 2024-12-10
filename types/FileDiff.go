package types

import (
	"os"
	"strings"
)

type FileDiff struct {
	OldFile     *FileChanges
	NewFile     os.DirEntry
	FullDirPath string
}

func (diff *FileDiff) IsValid() bool {
	//fmt.Println(diff.OldFile)
	if diff.OldFile == nil {
		return false
	}

	if diff.NewFile == nil {
		return false
	}

	return true
}

func (diff *FileDiff) IsDiff() bool {
	if !diff.IsValid() {
		return true
	}

	newInfo, newInfoError := diff.NewFile.Info()

	if newInfoError != nil {
		return true
	}

	if diff.OldFile.ModTime.Format("2006-01-02 15:04:05") != newInfo.ModTime().Format("2006-01-02 15:04:05") {
		return true
	}

	return false
}

func (diff *FileDiff) FullPath() string {
	return diff.FullDirPath + "\\" + diff.NewFile.Name()
}

func (diff *FileDiff) DirPaths(basePath string) (string, string) {
	//if diff.NewFile == nil {
	//	return false
	//}

	parsedDirPath := strings.Replace(diff.FullDirPath, basePath, "", 1)
	parsedFilePath := parsedDirPath + "\\" + diff.NewFile.Name()
	if parsedDirPath == "" {
		parsedDirPath = "\\"
	}

	return parsedDirPath, parsedFilePath
}

func (diff *FileDiff) GetCsvReg(progParams ProgParams) string {
	_, virtualFilePath := diff.DirPaths(progParams.BasePath)
	virtualFilePath = strings.ReplaceAll(virtualFilePath, "\\", "/")

	info, _ := diff.NewFile.Info()

	if diff.IsDiff() {
		return virtualFilePath + ";" + info.ModTime().Format("2006-01-02 15:04:05") + ";" + progParams.BackupDate
	}

	return ""
}
