package replication

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sergeyzalunin/go-replication-loader/logger"
)

// FileLoader contains all base behavior to work with replication files,
// its directories and description of replication
type FileLoader struct {
	log *logger.Log
	ReplicationDirectory string
}

// Init file loader
func (file *FileLoader) Init(dbName string, log *logger.Log) error {
	file.log = log
	if len(strings.TrimSpace(dbName)) == 0 {
		return errors.New("Database name didn't set in command line. Use -help command")
	}

	return setReplicationDirectory(file, dbName)
}

// GetFiles gets files from the replication
// directory by particular pattern: *.rep, *.desc, etc
func (file *FileLoader) GetFiles(pattern string) []string {
	pattern = filepath.Join(file.ReplicationDirectory, pattern)

	files, err := filepath.Glob(pattern)
	if err != nil {
		err := fmt.Errorf("An error occured while getting files from directory %s. %v", 
			file.ReplicationDirectory, err)
		file.log.Fatal(err)
		panic(err)
	}

	sort.Slice(files, func(i, j int) bool {
		fiBefore := getFileInfo(files[i], file.log)
		fiAfter := getFileInfo(files[j], file.log)
		if fiBefore != nil && fiAfter != nil {
			return fiBefore.ModTime().Before(fiAfter.ModTime())
		}
		return false
	})

	return files
}

func getFileInfo(file string, log *logger.Log) os.FileInfo {
	fi, err := os.Stat(file)
	if err != nil {
		err = fmt.Errorf("An error occured while getting file info %s. %v", file, err)
		log.Fatal(err)
		panic(err)
	}
	return fi
}

func setReplicationDirectory(file *FileLoader, dbName string) error {
	dir, err := createReplicationDirectory(dbName)
	if err != nil {
		return err
	}

	file.ReplicationDirectory = dir
	return nil
}

// createReplicationDirectory - creates a directory from which gets replication and description files
// Its name concats form Database name and "Replics" suffix
func createReplicationDirectory(dbName string) (string, error) {
	dir, err := getExecutableDirectory()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, dbName+"Replics")
	createDirectory(dir)
	return dir, nil
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
