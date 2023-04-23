// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error
	Log(tag string, message string)

	// Scheduler Services
	Schedule(timer time.Duration, h ServiceHandleFunc) chan bool /*exit channel*/
}

// Microservice is the centralized service management
type Microservice struct {
	echo        *echo.Echo
	exitChannel chan bool
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{
		echo: echo.New(),
	}
}

// startHTTP will start HTTP service, this function will block thread
func (ms *Microservice) startHTTP(exitChannel chan bool) error {
	// Caller can exit by sending value to exitChannel
	go func() {
		<-exitChannel
		ms.stopHTTP()
	}()
	return ms.echo.Start(":8080")
}

func (ms *Microservice) stopHTTP() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ms.echo.Shutdown(ctx)
}

// Start start all registered services
func (ms *Microservice) Start() error {

	httpN := len(ms.echo.Routes())
	var exitHTTP chan bool
	if httpN > 0 {
		exitHTTP = make(chan bool, 1)
		go func() {
			ms.startHTTP(exitHTTP)
		}()
	}

	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	ms.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		case <-ms.exitChannel:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		}
	}

	ms.Cleanup()
	return nil
}

// Stop stop the services
func (ms *Microservice) Stop() {
	if ms.exitChannel == nil {
		return
	}
	ms.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	ms.Log("MS", "Start cleanup")
	return nil
}

// Log log message to console
func (ms *Microservice) Log(tag string, message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println(tag+":", fns[len(fns)-1], line, message)
}
