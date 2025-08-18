package subsystems

import (
    "os"
    "fmt"
    "path"
    "strconv"
)

type CpuSubSystem struct {
}

func (s *CpuSubSystem) Name() string {
    return "cpu"
}

func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
    if res.CpuCfsQuota == 0 {
        return nil
    }
    subCgroupPath, err := getCgroupPath(cgroupPath, true)
    if err != nil {
        return err
    }
    // get string in form of "quota period" ("20000 100000" = 20%)
    quota := fmt.Sprintf("%s %s", strconv.Itoa(100000*res.CpuCfsQuota/100), 100000)
    if err = os.WriteFile(path.Join(subCgroupPath, "cpu.max"), []byte(quota), 0644); err != nil {
        return fmt.Errorf("error setting cgroup cpu share in cgroup[%s]: %v", cgroupPath, err)
    }
    return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int) error {
    return applyCgroup(pid, cgroupPath)
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
    subCgroupPath, err := getCgroupPath(cgroupPath, false)
    if err != nil {
        return err
    }
    return os.RemoveAll(subCgroupPath)
}
