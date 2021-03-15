package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/robfig/cron/v3"
)

func execute(command string, args []string) {

	fmt.Println("executing:", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	cmd.Wait()
}

func create() (cr *cron.Cron, wgr *sync.WaitGroup) {
	var schedule string = os.Args[1]
	var command string = os.Args[2]
	var args []string = os.Args[3:len(os.Args)]

	wg := &sync.WaitGroup{}

	c := cron.New()
	fmt.Println("new cron:", schedule)

	_, err := c.AddFunc(schedule, func() {
		wg.Add(1)
		execute(command, args)
		wg.Done()
	})

	if err != nil {
		fmt.Println(err)
		stop(c, wg)
	}
	return c, wg
}

func start(c *cron.Cron, wg *sync.WaitGroup) {
	c.Start()
}

func stop(c *cron.Cron, wg *sync.WaitGroup) {
	fmt.Println("Stopping")
	c.Stop()
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Exiting")
	os.Exit(0)
}

func main() {

	c, wg := create()

	go start(c, wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	stop(c, wg)
}
