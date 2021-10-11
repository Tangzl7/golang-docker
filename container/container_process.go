package container

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"syscall"
	"os/exec"
	// "path/filepath"
	log "github.com/Sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	// args := append([]string{"init"}, cmdArray...)
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | 
					syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Dir ="/root/busybox"
	// add a file descriptor for process
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("get user command error")
	}

	setUpMount()
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

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func readUserCommand() []string {
	// uintptr(3) means the file descriptor of index 3
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func pivotRoot(root string) error {
	// pivot_root to new root rootfs, move old_root to rootfs/.pivot_root
	if err := syscall.Chroot(root); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// make current path be root path
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}
	if err := os.RemoveAll("proc"); err != nil {
		return fmt.Errorf("failed to remove rootfs/proc: %v", err)
	}
	if err := os.Mkdir("proc", 0755); err != nil {
		return fmt.Errorf("failed to mkdir rootfs/proc: %v", err)
	}
	// remove pivot_root
	return nil
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	pivotRoot(pwd)

	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=775")
}