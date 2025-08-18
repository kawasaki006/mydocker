package subsystems

import (
    "os"
    "path"
    "fmt"
    "strconv"
)

func getCgroupPath(cgroupPath string, autoCreate bool) (string, error) {
    cgroupRoot := "/sys/fs/cgroup"
    fullPath := path.Join(cgroupRoot, cgroupPath)
    if !autoCreate {
        return fullPath, nil
    }
    // create if directory doesn't exist
    if _, err := os.Stat(fullPath); err != nil {
        if os.IsNotExist(err) {
            if err := os.Mkdir(fullPath, 0755); err != nil {
                return fullPath, fmt.Errorf("error creating cgroup directory: %v", err)
            }
        } else {
                return fullPath, fmt.Errorf("error stat cgroup path: %v", err)
        }
    }
    return fullPath, nil
}

func applyCgroup(pid int, cgroupPath string) error {
    subCgroupPath, err := getCgroupPath(cgroupPath, true)
    if err != nil {
        return fmt.Errorf("error getting cgroup path when adding process: %v", err)
    }
    // write pid to cgroup.procs file
    if err = os.WriteFile(path.Join(subCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
        return fmt.Errorf("error adding process to cgroup[%s]: %v", subCgroupPath, err)
    }
    return nil
}
