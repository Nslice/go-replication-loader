package loader

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/logger"
)

// ProcessExecutor is a struct that allows to run processes
type ProcessExecutor struct {
	log                      logger.Log
	PathToAdminToolsConsole  string
	PathToCompilationPluting string
	dir                      string
}

// NewProcessExecutor is a ProcessExecutor factory
func NewProcessExecutor(args argsp.ArgumentOptions, log logger.Log) ProcessExecutor {
	return ProcessExecutor{
		log,
		filepath.Join(args.WorkingDirectory, "Akforta.eLeed.AdminToolsConsole.exe"),
		filepath.Join(args.WorkingDirectory, "BIZ.Compiler.exe"),
		args.WorkingDirectory,
	}
}

// RunAdminToolsConsole starts the installation process
func (p *ProcessExecutor) RunAdminToolsConsole(args string) {
	p.run(p.PathToAdminToolsConsole, args)
}

// RunCompilationPluting starts the compilation process
func (p *ProcessExecutor) RunCompilationPluting(args string) {
	p.run(p.PathToCompilationPluting, args)
}

func (p *ProcessExecutor) run(filename string, args string) {
	if _, err := exec.LookPath(filename); err == nil {
		cmd := exec.Command(filename)

		// Filename + args sets directly due to avoid auto arguments escaping.
		// Akforta.eLeed.AdminToolsConsole.exe can't handle escaped arguments
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
			CmdLine:    filename + args,
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			p.logProcess(output, cmd, err)
			panic(err)
		}
		p.logProcess(output, cmd, nil)
		p.eleedSpecificCheckings(filename, cmd.ProcessState.ExitCode())
	} else {
		p.logProcess([]byte{}, nil, err)
		panic(err)
	}
}

func (p *ProcessExecutor) logProcess(output []byte, cmd *exec.Cmd, err error) {
	if err == nil {
		p.log.Info(string(output))
	} else {
		p.log.Error(err, string(output))
	}

	if cmd != nil {
		p.log.Info("exit code: ", cmd.ProcessState.ExitCode())
	}
}

func (p *ProcessExecutor) eleedSpecificCheckings(filename string, exitCode int) {
	msg := fmt.Sprintf("The execution of the programm %s was completed with code %d", filename, exitCode)
	if filename == p.PathToCompilationPluting && exitCode != 0 {
		p.log.Error(msg)
	} else {
		p.log.Info(msg)
	}
}
