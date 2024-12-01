package main

import (
	"os"
	"strings"
)

type FileDiff struct {
	OldFile, NewFile os.DirEntry
	FullDirPath      string
}

func (diff *FileDiff) IsValid() bool {
	if diff.OldFile == nil {
		return false
	}

	if diff.NewFile == nil {
		return false
	}

	//if not new file => panic

	return true
}

func (diff *FileDiff) IsDiff() bool {
	if !diff.IsValid() {
		return true
	}

	oldInfo, oldInfoError := diff.OldFile.Info()
	newInfo, newInfoError := diff.NewFile.Info()

	if oldInfoError != nil || newInfoError != nil {
		return true
	}

	if oldInfo.ModTime() != newInfo.ModTime() {
		return true
	}

	return false
}

func (diff *FileDiff) DirPaths(basePath string) (string, string) {
	//if diff.NewFile == nil {
	//	return false
	//}

	parsedDirPath := strings.Replace(diff.FullDirPath, basePath, "", 1)
	parsedFilePath := parsedDirPath + "\\" + diff.NewFile.Name()

	return parsedDirPath, parsedFilePath
}
