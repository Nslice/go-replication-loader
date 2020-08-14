package services

import (
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
)

var (
	states = map[svc.State]string{
		svc.State(windows.SERVICE_STOPPED):          "Stopped",
		svc.State(windows.SERVICE_START_PENDING):    "StartPending",
		svc.State(windows.SERVICE_STOP_PENDING):     "StopPending",
		svc.State(windows.SERVICE_RUNNING):          "Running",
		svc.State(windows.SERVICE_CONTINUE_PENDING): "ContinuePending",
		svc.State(windows.SERVICE_PAUSE_PENDING):    "PausePending",
		svc.State(windows.SERVICE_PAUSED):           "Paused",
	}
)
