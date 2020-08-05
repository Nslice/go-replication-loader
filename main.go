package main

// LineBreak is the windows line break
const LineBreak = "\r\n"

func main() {
	args := new(ArgumentOptions)
	args.Init()
}

func getServiceWorker(args ArgumentOptions) IService{
	service := new(ServiceWorker)
	service.ServiceName = args.ConsoleServiceName
	return service
}