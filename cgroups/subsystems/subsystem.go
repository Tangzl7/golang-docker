package subsystems

// memory limit, cpu time slice weight, cpu core number
type ResourceConfig struct {
	MemoryLimit string
	CpuShare string
	CpuSet string
}

type Subsystem interface {
	// subsystem name, e.g cpu memory
	Name() string
	// set cgroup's resource limit
	Set(path string, res *ResourceConfig) error
	// add process to a cgroup
	Apply(path string, pid int) error
	// remove a cgroup
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem {
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)