package container

import (
    "os"
    "os/exec"
    "syscall"

    log "github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
    return cmd, writePipe
}
