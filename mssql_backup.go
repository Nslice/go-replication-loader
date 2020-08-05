package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/denisenkom/go-mssqldb"
)

// DoBackup create an backup of target database provided via ArgumentOptions
func DoBackup(args ArgumentOptions) error  {
	connString, err := getConnection(args, 7200)
	if err != nil {
		return err
	}

	connector, err := mssql.NewConnector(connString)
    if err != nil {
        return err
    }

	db := sql.OpenDB(connector)
	defer db.Close()

	backupCommand := getBackupCommand(args)
	
	_, err = db.Exec(backupCommand)	
	return err
}

func getConnection(args ArgumentOptions, timeout int) (string, error) {	
	if args.DatabaseName == "" {
		err := "The database name doesn't set in command line. Use --dbname or --help command"
		return "", fmt.Errorf(err)
	}
	
	connectionString := fmt.Sprintf("Server=%s; Database=%s; Trusted_Connection=true; Connection Timeout=%d", 
		args.DbDataSource, args.DatabaseName, timeout);

	return connectionString, nil
}

func getBackupCommand(args ArgumentOptions) string {
	log.Printf("Backup of database %s - start", args.DatabaseName)

	filename := filepath.Join(args.BackupPath, args.DatabaseName + "_ReplicLoaderAutobackup.bak")
	
	var compressionString string
	if args.UseCompression {
		compressionString = "COMPRESSION,"
	}

	result := fmt.Sprintf("BACKUP DATABASE %s TO DISK = %s ", args.DatabaseName, filename)
	result += fmt.Sprintf("INIT, NAME = N'%s Database Backup',")
	result += fmt.Sprintf("SKIP, NOREWIND, NOUNLOAD, %s STATS = 10", compressionString)

	return result
}