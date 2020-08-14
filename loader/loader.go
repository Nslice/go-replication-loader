package loader

import (
	"fmt"
	"os"
	"strings"

	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/logger"
	"github.com/sergeyzalunin/go-replication-loader/mssql"
	"github.com/sergeyzalunin/go-replication-loader/replication"
	"github.com/sergeyzalunin/go-replication-loader/services"
)

// Loader is a container that has instances to load replications
type Loader struct {
	log            logger.Log
	args           argsp.ArgumentOptions
	repl           replication.ReplicationLoader
	consoleService services.IService
	netpipeService services.IService
	executor       ProcessExecutor
}

// NewLoader is a constructor to create a new Loader struct
func NewLoader(args argsp.ArgumentOptions, log logger.Log) *Loader {
	repl := replication.ReplicationLoader{}
	repl.Init(args)

	console := services.NewService(args.ConsoleServiceName, log)
	netpipe := services.NewService(args.NetPipeServiceName, log)
	executor := NewProcessExecutor(args, log)

	return &Loader{log, args, repl, console, netpipe, executor}
}

// Load starts the process of loading
func (l *Loader) Load() error {
	replicationFiles := l.repl.GetReplicationFiles()

	if len(replicationFiles) > 0 {
		l.preloadingProcess()

		for _, rep := range replicationFiles {
			l.log.Info("The replication ", rep, " is loading")

			args := l.getAdminToolsConsoleArguments(rep)
			l.executor.RunAdminToolsConsole(args)
			err := os.Remove(rep)
			if err == nil {
				l.log.Info("The replication file ", rep, " was deleted from folder")
			} else {
				msg := fmt.Sprintf("Failed to remove a replication file %s\n", rep)
				l.log.Error(msg, err)
			}
		}

		l.postloadingProcesses()
	}
	return nil
}

func (l *Loader) preloadingProcess() {
	l.log.Info("Replication(s) is in the directory ", l.repl.ReplicationDirectory)

	err := l.netpipeService.StopService()
	l.log.LogIfError(err, "Failed stop the netpipe service")

	err = l.consoleService.StopService()
	if err != nil {
		msg := "Failed to stop the console monolithic service"
		l.log.Fatal(msg)
		panic(msg)
	}

	mssql.DoBackup(l.args, l.log)
	err = l.consoleService.StartService()
	if err != nil {
		msg := "Failed to start the console monolithic service\n"
		l.log.Fatal(msg, err)
		panic(msg)
	}
}

func (l *Loader) postloadingProcesses() {
	args := l.getCompilationPluginArguments()
	l.executor.RunCompilationPluting(args)

	l.netpipeService.StartService()
	l.log.Info("All replications have already loaded successfully")
}

func (l *Loader) getCompilationPluginArguments() string {
	result := strings.Builder{}
	result.Grow(100)
	result.WriteString(getArgument("user", l.args.User))
	result.WriteString(getArgument("password", l.args.Password))
	return result.String()
}

func (l *Loader) getAdminToolsConsoleArguments(rep string) string {
	result := strings.Builder{}
	result.Grow(200)
	result.WriteString(getArgument("plugin", "InnerReplicationPlugin"))
	result.WriteString(getArgument("user", l.args.User))
	result.WriteString(getArgument("password", l.args.Password))
	result.WriteString(" --import")
	result.WriteString(" --nocompilation")
	result.WriteString(getIntArgument("verbose", 4))
	result.WriteString(getArgument("file", rep))

	return result.String()
}

func getArgument(key, value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf(" --%s \"%s\"", key, value)
}

func getIntArgument(key string, value int) string {
	return fmt.Sprintf(" --%s %d", key, value)
}
