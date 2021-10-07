package subsystems

import (
	"os"
	"fmt"
	"path"
	"bufio"
	"strings"
)

// get subsystem mount point by /proc/self/mountinfo
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields) - 1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// get cgroup's absolut path
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgrouptRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgrouptRoot, cgroupPath)); err == nil || 
		(autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgrouptRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgrouptRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}