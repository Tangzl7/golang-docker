package main

import (
	"fmt"
	"golang-docker/container"
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
)

var runCommand = &cli.Command {
	Name: "run",
	Usage: "Create a container with namespace and cgroup golang-docker run -ti [command]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
	},
	// run action
	Action: func(context *cli.Context) error {
		if context.NArg() < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = &cli.Command {
	Name: "init",
	Usage: "Init container process run user's process in container. Do not call it out side",
	// init action
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}