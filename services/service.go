// +build windows

package services

import (
	"fmt"
	"time"

	"github.com/sergeyzalunin/go-replication-loader/logger"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type serviceFunc func(*mgr.Service) error

// IService - base functions for working with services
type IService interface {
	HasService() error
	StartService() error
	StopService() error
}

// ServiceWorker is a handler to work with windows services
type ServiceWorker struct {
	log         *logger.Log
	ServiceName string
}

// NewService is a constructor to get IService
func NewService(serviceName string, log *logger.Log) IService {
	return ServiceWorker{log, serviceName}
}

// HasService returns true
// if window already has the service
// with name presented in ServiceWorker struct
func (worker ServiceWorker) HasService() error {
	return worker.serviceAction(worker.hasService)
}

func (worker ServiceWorker) hasService(service *mgr.Service) error {
	if service != nil {
		return nil
	}
	return fmt.Errorf("service %s does not exist", worker.ServiceName)
}

// StartService starts the service
// with name from ServiceWorker struct
func (worker ServiceWorker) StartService() error {
	return worker.serviceAction(worker.startService)
}

func (worker ServiceWorker) startService(service *mgr.Service) error {
	if worker.hasServiceStatus(service, svc.Stopped) {
		err := service.Start()
		if err != nil {
			return fmt.Errorf("Could not start the service: %v", err)
		}

		status, err := worker.waitingForState(service, svc.Running)
		if err != nil {
			return fmt.Errorf("Failed to stop service %s, status: %v", worker.ServiceName, states[status.State])
		}

		worker.log.Info("Service ", worker.ServiceName, " is already started, final state: ", states[status.State])
	}

	return nil
}

// StopService stops the service
// with name from ServiceWorker struct
func (worker ServiceWorker) StopService() error {
	return worker.serviceAction(worker.stopService)
}

func (worker ServiceWorker) stopService(service *mgr.Service) error {
	if worker.hasServiceStatus(service, svc.Running) {
		status, err := service.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("Could not stop the service: %v", err)
		}

		status, err = worker.waitingForState(service, svc.Stopped)
		if err != nil {
			return fmt.Errorf("Failed to stop service %s, status: %v", worker.ServiceName, states[status.State])
		}

		worker.log.Info("Service ", worker.ServiceName, "is already stopped, final state: ", states[status.State])
	}

	return nil
}

func (worker ServiceWorker) serviceAction(action serviceFunc) error {
	if worker.ServiceName == "" {
		return nil
	}

	manager, err := mgr.Connect()
	if err != nil {
		err = fmt.Errorf("Cannot connect to manager %v", err)
		worker.log.Error(err)
		return err
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(worker.ServiceName)
	if err != nil {
		err = fmt.Errorf("Service '%s' does not exist: %v", worker.ServiceName, err)
		worker.log.Error(err)
		return err
	}
	defer service.Close()

	err = action(service)
	if err != nil {
		worker.log.Error(err)
	}

	return err
}

func (worker ServiceWorker) hasServiceStatus(service *mgr.Service, state svc.State) bool {
	status, err := service.Query()
	if err != nil {
		errMsg := fmt.Errorf("Failed to get service '%s' status: %v", worker.ServiceName, err)
		worker.log.Error(errMsg)
		return false
	}
	worker.log.Info("Service ", worker.ServiceName, " state is ", states[status.State])
	return status.State == state
}

func (worker ServiceWorker) waitingForState(service *mgr.Service, state svc.State) (svc.Status, error) {
	start := time.Now()

	var status svc.Status
	var err error

	for 600*time.Second > time.Since(start) {
		status, err = service.Query()
		if err != nil {
			err = fmt.Errorf("failed to get service '%s' status: %v", worker.ServiceName, err)
			worker.log.Error(err)
			break
		}

		if status.State == state {
			break
		}
	}

	return status, err
}
