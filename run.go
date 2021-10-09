package main

import (
	"os"
	"strings"
	"golang-docker/container"
	"golang-docker/cgroups"
	"golang-docker/cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res* subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// send cmds to son process by writing cmds in pipe
	sendInitCommand(cmdArray, writePipe)
	cgroupManager := cgroups.NewCgroupManager("golang-docker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}