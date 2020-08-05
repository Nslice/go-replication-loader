package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileLoader contains all base behavior to work with replication files,
// its directories and description of replication
type FileLoader struct {
	ReplicationDirectory string
}

// Init file loader
func (file *FileLoader) Init(args ArgumentOptions) error {	
	dbName := args.DatabaseName
	if len(strings.TrimSpace(dbName)) == 0 {
		return errors.New("Database name didn't set in command line. Use -help command")
	}

	setReplicationDirectory(file, dbName)

	return nil
}

// GetFiles gets files from the replication 
// directory by particular pattern: *.rep, *.desc, etc
func (file *FileLoader) GetFiles(pattern string) []string {
	pattern = filepath.Join(file.ReplicationDirectory, pattern)

	files, err := filepath.Glob(file.ReplicationDirectory)	
	if err != nil {
		log.Fatalf("An error occured while getting files from directory %v. %v", pattern, err)
	}

	sort.Slice(files, func(i, j int) bool{
		fiBefore := getFileInfo(files[i])
		fiAfter := getFileInfo(files[j])
		if fiBefore != nil && fiAfter != nil {
			return fiBefore.ModTime().Before(fiAfter.ModTime())
		}
		return false
	})

	return files
}

func getFileInfo(file string) os.FileInfo {
	fi, err := os.Stat(file)
	if err != nil {
		log.Fatalf("An error occured while getting file info %v. %v", file, err)
		return nil
	}
	return fi
}

func setReplicationDirectory(file *FileLoader, dbName string) error {
	dir, err := getExecutableDirectory()
	if err != nil {
		return err
	}

	file.ReplicationDirectory = filepath.Join(dir, dbName + "Replics")
	createDirectory(file.ReplicationDirectory)
	return nil
}

func getExecutableDirectory() (string, error) {
	ex, err := os.Executable()
    if err != nil {
        return "", err
    }
	exPath := filepath.Dir(ex)
	return exPath, nil
}

func createDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}