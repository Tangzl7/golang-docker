package main

import (
	"fmt"
	"golang-docker/container"
	"golang-docker/cgroups/subsystems"
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

var runCommand = &cli.Command {
	Name: "run",
	Usage: "Create a container with namespace and cgroup golang-docker run -ti [command]",
	Flags: []cli.Flag {
		&cli.BoolFlag {
			Name: "ti",
			Usage: "enable tty",
		},
		&cli.StringFlag {
			Name: "m",
			Usage: "memory limit",
		},
		&cli.StringFlag {
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag {
			Name: "cpuset",
			Usage: "cpuset limit",
		},
	},
	// run action
	Action: func(context *cli.Context) error {
		if context.NArg() < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmdArray := context.Args().Slice()
		tty := context.Bool("ti")
		resConf := &subsystems.ResourceConfig {
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpushare"),
		}

		Run(tty, cmdArray, resConf)
		return nil
	},
}

var initCommand = &cli.Command {
	Name: "init",
	Usage: "Init container process run user's process in container. Do not call it out side",
	// init action
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmdArray := context.Args().Slice()
		err := container.RunContainerInitProcess(cmdArray)
		return err
	},
}