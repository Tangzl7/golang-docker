package main

import (
	"os"
	"golang-docker/container"
	"golang-docker/cgroups"
	"golang-docker/cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res* subsystems.ResourceConfig) {
	parent := container.NewParentProcess(tty, cmdArray)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("golang-docker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	parent.Wait()
	os.Exit(-1)
}