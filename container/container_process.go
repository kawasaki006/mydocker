package container

import (
    "os"
    "os/exec"
    "syscall"
    "path"

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
    } else {
        // redirect stdout and stderr to log file
        dirPath := path.Join(InfoLocation, containerId)
        if err = os.MkdirAll(dirPath, 0622); err != nil {
            log.Errorf("error creating info dir in newparentproc [%s]: %v", dirPath, err)
            return nil, nil
        }
        logFilePath := dirPath + GetLogFile(containerId)
        logFile, err := os.Create(logFilePath)
        if err != nil {
            log.Errorf("error creating log file in newparentproc [%s]: %v", logFilePath, err)
            return nil, nil
        }
        cmd.Stdout = logFile
        cmd.Stderr = logFile
    }
    // carry read pipe as extra file
    cmd.ExtraFiles = []*os.File{readPipe}
    // mount overlay fs
    NewWorkspace(containerId, imageName, volume)
    // set mount point as container working dir
    cmd.Dir = utils.GetMerged(containerId)
    return cmd, writePipe
}
