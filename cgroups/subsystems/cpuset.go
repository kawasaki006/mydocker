package subsystems

import (
    "os"
    "fmt"
    "path"
)

type CpusetSubSystem struct {
}

func (s *CpusetSubSystem) Name() string {
    return "cpuset"
}

func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
    if res.CpuSet == "" {
        return nil
    }
    subCgroupPath, err := getCgroupPath(cgroupPath, true)
    if err != nil {
        return err
    }
    if err := os.WriteFile(path.Join(subCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
        return fmt.Errorf("error setting cpuset in cgroup[%s]: %v", cgroupPath, err)
    }
    return nil
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
    return applyCgroup(pid, cgroupPath)
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
    subCgroupPath, err := getCgroupPath(cgroupPath, false)
    if err != nil {
        return err
    }
    return os.RemoveAll(subCgroupPath)
}
