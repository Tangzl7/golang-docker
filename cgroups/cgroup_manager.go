package cgroups

import (
	"golang-docker/cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup path in hierarchy
	Path string
	// resource config
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// add pid into cgroup
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// add resource limit for cgroup
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// remove cgroup for all subsystems
func (c *CgroupManager) Destory() error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
