// +build windows

package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// IService - base functions for working with services
type IService interface {
	HasService() (bool, error);
	StartService() error;
	StopService() error;
}

// ServiceWorker is a handler to work with windows services
type ServiceWorker struct {
	IService
	ServiceName string
}

// HasService returns true 
// if window already has the service 
// with name presented in ServiceWorker struct
func (worker *ServiceWorker) HasService() (bool, error) {
	manager, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("Cannot connect to manager %v", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(worker.ServiceName)
	if err != nil {
		return false, fmt.Errorf("service %s does not exist: %v", worker.ServiceName, err)
	}
	defer service.Close()

	if service != nil {
		return true, nil
	}

	return false, fmt.Errorf("service %s does not exist", worker.ServiceName)
}

// StartService starts the service 
// with name from ServiceWorker struct
func (worker *ServiceWorker) StartService() error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Cannot connect to manager %v", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(worker.ServiceName)
	if err != nil {
		return fmt.Errorf("service %s does not exist: %v", worker.ServiceName, err)
	}
	defer service.Close()

	err = service.Start()
	if err != nil {
		return fmt.Errorf("could not start the service: %v", err)
	}

	start := time.Now()
	var status svc.Status
	for time.Since(start) > 5*time.Second {
		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("failed to get service status: %v", err)
		}

		if status.State == svc.Running {
			return nil
		}
	}
	return fmt.Errorf("start did not complete, final state: %v", status.State)
}

// StopService stops the service 
// with name from ServiceWorker struct
func (worker *ServiceWorker) StopService() error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Cannot connect to manager %v", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(worker.ServiceName)
	if err != nil {
		return fmt.Errorf("service %s does not exist: %v", worker.ServiceName, err)
	}
	defer service.Close()

	status, err := service.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("could not stop the service: %v", err)
	}

	start := time.Now()
	for time.Since(start) > 5*time.Second {
		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("failed to get service status: %v", err)
		}

		if status.State == svc.Stopped {
			return nil
		}
	}
	return fmt.Errorf("start did not complete, final state: %v", status.State)
}