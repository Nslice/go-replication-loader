package loader

import (
	"github.com/sergeyzalunin/go-replication-loader/logger"
	"github.com/sergeyzalunin/go-replication-loader/services"
)

// Builder serves to construct a process of installing replications
type Builder struct {
	loader Loader
	log    *logger.Log
}

// NewLoaderBuilder is a constructor for Builder
func NewLoaderBuilder(log *logger.Log) *Builder {
	return &Builder{Loader{}, log}
}

// AddCMS allows to get console monolithic service by its name
func (b *Builder) AddCMS(serviceName string) {
	service, err := getService(serviceName, b.log)
	if err != nil {
		panic(err)
	}
	b.loader.consoleService = service
}

// AddNetPipe allows to get netpipe service by its name
func (b *Builder) AddNetPipe(serviceName string) {
	if serviceName == "" {
		return
	}

	service, err := getService(serviceName, b.log)
	if err != nil {
		panic(err)
	}

	b.loader.netpipeService = service
}

// Build the replication loader
func (b *Builder) Build() Loader {
	b.log.Info("Replication(s) ", b, " is in the directory")
	return b.loader
}

func getService(serviceName string, log *logger.Log) (services.IService, error) {
	service := services.NewService(serviceName, log)
	if err := service.HasService(); err != nil {
		log.Error(err)
		return service, err
	}
	return service, nil
}
