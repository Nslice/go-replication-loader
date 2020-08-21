package mssql

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/go-errors/errors"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/logger"
)

// DoBackup create an backup of target database provided via ArgumentOptions
func DoBackup(args *argsp.ArgumentOptions, log *logger.Log) {
	if args.SkipBackup {
		log.Info("Backup of database skipped due to skipbackup flag")
		return
	}

	backupState := "with error"
	log.Info("Backup of database ", args.DatabaseName, " started")
	defer func() {
		log.Info("Backup of database ", args.DatabaseName, " finished ", backupState)
	}()

	err := doBackup(args, log)
	if err == nil {
		backupState = "successfully"
	} else {
		msg := "You are not allowed to continue installation a replication without successful backup"
		log.Fatal(errors.Errorf(msg))
		panic(msg)
	}
}

func doBackup(args *argsp.ArgumentOptions, log *logger.Log) error {
	connString, err := getConnection(args)
	if err != nil {
		return err
	}

	backupCommand := getBackupCommand(args, log)

	connector, err := mssql.NewConnector(connString)
	if err != nil {
		return err
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	_, err = db.Exec(backupCommand)
	return err
}

func getConnection(args *argsp.ArgumentOptions) (string, error) {
	if args.DatabaseName == "" {
		err := "The database name doesn't set in command line. Use -dbname or -help command"
		return "", fmt.Errorf(err)
	}
	connectionString := NewConnectionString(args)
	return connectionString, nil
}

func getBackupCommand(args *argsp.ArgumentOptions, log *logger.Log) string {
	filename := filepath.Join(args.BackupPath, args.DatabaseName+"_ReplicLoaderAutobackup.bak")
	log.Info("Backup will be saved at the path ", filename)

	var compressionString string
	if args.UseCompression {
		compressionString = "COMPRESSION,"
	}

	result := fmt.Sprintf("BACKUP DATABASE %s TO DISK = '%s' ", args.DatabaseName, filename)
	result += fmt.Sprintf("WITH NOFORMAT, INIT, NAME = N'%s Database Backup', ", args.DatabaseName)
	result += fmt.Sprintf("SKIP, NOREWIND, NOUNLOAD, %s STATS = 10", compressionString)
	log.Info("Backup sql query: ", result)

	return result
}
