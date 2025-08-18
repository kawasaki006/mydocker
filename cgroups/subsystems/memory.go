package subsystems

import (
    "os"
    "path"
    "fmt"
)

type MemorySubSystem struct {
}

func (s *MemorySubSystem) Name() string {
    return "memory"
}

// set memory limit (specified by res) to cgroup
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
    if res.MemoryLimit == "" {
        return nil
    }
    subCgroupPath, err := getCgroupPath(cgroupPath, true)
    if err != nil {
        return err
    }
    // set memory limit
    if err := os.WriteFile(path.Join(subCgroupPath, "memory.max"), []byte(res.MemoryLimit), 0644); err != nil {
        return fmt.Errorf("error setting cgroup memory: %v", err)
    }
    return nil
}

// add process to cgroup
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
    return applyCgroup(pid, cgroupPath)
}

// remove cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
    subCgroupPath, err := getCgroupPath(cgroupPath, false)
    if err != nil {
        return fmt.Errorf("error getting cgroup path in cgroup removal: %v", err)
    }
    return os.RemoveAll(subCgroupPath)
}
