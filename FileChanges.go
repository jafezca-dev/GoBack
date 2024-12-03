package main

import "time"

type FileChanges struct {
	ModTime   time.Time
	BackupTag string
}
