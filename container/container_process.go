package container

import (
	"os"
	"fmt"
	"syscall"
	"os/exec"
	log "github.com/Sirupsen/logrus"
)

func NewParentProcess(tty bool, cmdArray []string) *exec.Cmd {
	args := append([]string{"init"}, cmdArray...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | 
					syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func RunContainerInitProcess(cmdArray []string) error {
	log.Infof("command %s", cmdArray)
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("get user command error")
	}

	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}