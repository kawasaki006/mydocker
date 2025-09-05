package container

import (
    "os"
    "os/exec"
    "syscall"

    log "github.com/sirupsen/logrus"
    "github.com/kawasaki006/mydocker/utils"
)

func NewParentProcess(tty bool, volume, containerId, imageName string) (*exec.Cmd, *os.File) {
    // initialize pipe: read pipe for init proc to read command string; write pipe for run proc to send command
    readPipe, writePipe, err := os.Pipe()
    if err != nil {
        log.Errorf("Error initializing new pipe: %v", err)
        return nil, nil
    }

    cmd := exec.Command("/proc/self/exe", "init")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
    }
    if (tty) {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }
    // carry read pipe as extra file
    cmd.ExtraFiles = []*os.File{readPipe}
    // mount overlay fs
    NewWorkspace(containerId, imageName, volume)
    // set mount point as container working dir
    cmd.Dir = utils.GetMerged(containerId)
    return cmd, writePipe
}
