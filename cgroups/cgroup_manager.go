package cgroups

import (
    log "github.com/sirupsen/logrus"
    "github.com/kawasaki006/mydocker/cgroups/subsystems"
)

// path of the cgroup
type CgroupManager struct {
    Path string
    Resource *subsystems.ResourceConfig
    Subsystems []subsystems.Subsystem
}

// new a manager
func NewCgroupManager(path string) *CgroupManager {
    return &CgroupManager{
        Path:       path,
        Subsystems: subsystems.SubsystemsIns,
    }
}

// add process to all subsystems of the cgroup of this manager
func (c *CgroupManager) Apply(pid int) error {
    for _, subSysIns := range c.Subsystems {
        if err := subSysIns.Apply(c.Path, pid); err != nil {
            log.Errorf("error adding process[pid:%n] to cgroup subsystem[%s]: %v", pid, subSysIns.Name(), err)
        }
    }
    return nil
}

// set all subsystems of a manager
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
    for _, subSysIns := range c.Subsystems {
        if err := subSysIns.Set(c.Path, res); err != nil {
            log.Errorf("error setting subsystem[%s] of cgroup[%s]: %v", subSysIns.Name(), c.Path, err)
        }
    }
    return nil
}

// remove all subsystems of a manager
func (c *CgroupManager) Destroy() error {
    log.Infof("cleaning up cgroup....")
    for _, subSysIns := range c.Subsystems {
        if err := subSysIns.Remove(c.Path); err != nil {
            log.Errorf("error removing cgroup[%s]: %v", c.Path, err)
        }
    }
    return nil
}
