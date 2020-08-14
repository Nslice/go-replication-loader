package main

import (
	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/loader"
	"github.com/sergeyzalunin/go-replication-loader/logger"
	"github.com/sergeyzalunin/go-replication-loader/message"
)

func main() {
	var args argsp.ArgumentOptions
	var log logger.Log

	defer func() {
		sendEmail(args, log, recover())
	}()

	args = getArguments()

	log = logger.NewLogger(args.ProjectName)
	defer log.Close()

	err := doInstallation(args, log)
	if err != nil {
		panic(err)
	}
}

func getArguments() argsp.ArgumentOptions {
	args := argsp.ArgumentOptions{}
	args.Init()
	return args
}

func doInstallation(args argsp.ArgumentOptions, log logger.Log) error {
	l := loader.NewLoader(args, log)
	err := l.Load()
	return err
}
func sendEmail(args argsp.ArgumentOptions, log logger.Log, err interface{}) {
	e := message.New(args, log)
	if err == nil {
		e.Send()
	} else {
		e.SendFailed(err)
	}
}
