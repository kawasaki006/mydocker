package main

import (
    "os"
    "strings"

    "github.com/kawasaki006/mydocker/container"
    "github.com/kawasaki006/mydocker/cgroups"
    "github.com/kawasaki006/mydocker/cgroups/subsystems"
    log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume, containerName, imageName string) {
    // TODO: random id generation
    containerId := container.GenerateContainerId()
    
    parent, writePipe := container.NewParentProcess(tty, volume, containerId, imageName)
    // defer overlayfs cleanup
    // TODO: volume = " " for now
    if parent == nil {
        log.Errorf("Init process eror")
        return
    }

    // start container proc
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    // set up cgroup and add new contaier proc in
    cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
    cgroupManager.Set(res)
    cgroupManager.Apply(parent.Process.Pid)

    // send client full command
    sendInitCommand(cmdArray, writePipe)
    // close pipe
    writePipe.Close()

    // record container info
    _, err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerId, volume)
    if err != nil {
        log.Errorf("error recording container info: %v", err)
        return
    }
    
    if tty {
        parent.Wait()
        container.DeleteWorkspace(containerId, volume)
        container.DeleteContainerInfo(containerId)
        cgroupManager.Destroy()
    }

    go func() {
        if !tty {
            _, _ = parent.Process.Wait()
        }
        container.DeleteWorkspace(containerId, volume)
        container.DeleteContainerInfo(containerId)
        cgroupManager.Destroy()
    }()
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
    command := strings.Join(cmdArray, " ")
    log.Infof("Full command: %s", command)
    writePipe.WriteString(command)
}

