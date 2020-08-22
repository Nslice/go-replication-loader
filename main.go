package main

import (
	"encoding/json"
	"fmt"

	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/loader"
	"github.com/sergeyzalunin/go-replication-loader/logger"
	"github.com/sergeyzalunin/go-replication-loader/message"
)

func main() {
	var hasreplications bool
	var args *argsp.ArgumentOptions
	var log *logger.Log

	defer func() {
		if hasreplications {
			sendEmail(args, log, recover())
		}
	}()

	args = getArguments(log)

	log = logger.NewLogger(args.ProjectName)
	defer log.Close()

	hasreplications, err := doInstallation(args, log)
	if err != nil {
		panic(err)
	}
}

func getArguments(log *logger.Log) *argsp.ArgumentOptions {
	args := &argsp.ArgumentOptions{}
	args.Init()

	prjname, saveArgs, readSavedArgs := args.ProjectName, args.SaveArgs, args.ReadSavedArgs
	interactive, skipBackup := args.UseInteractive, args.SkipBackup

	savedArguments := argsp.ReadArguments(log)
	if !savedArguments.IsEmpty() {
		args = savedArguments
		args.UseInteractive = interactive
		args.ProjectName = prjname
		args.SaveArgs = saveArgs
		args.ReadSavedArgs = readSavedArgs
		args.SkipBackup = skipBackup
	}
	
	args = argsp.StartInteractiveMode(args, log)
	argsp.SaveArguments(args, log)

	if readSavedArgs {
		prettyPrint(args)
	}

	return args
}

func prettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func doInstallation(args *argsp.ArgumentOptions, log *logger.Log) (bool, error) {
	l := loader.NewLoader(args, log)
	return l.Load()
}
func sendEmail(args *argsp.ArgumentOptions, log *logger.Log, err interface{}) {
	e := message.New(args, log)
	if err == nil {
		e.Send()
	} else {
		e.SendFailed(err)
	}
}
