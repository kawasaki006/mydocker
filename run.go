package main

import (
    "os"
    "strings"

    "github.com/kawasaki006/mydocker/container"
    log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmdArray []string) {
    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
        log.Errorf("Init process eror")
        return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }
    // send client full command
    sendInitCommand(cmdArray, writePipe)
    // close pipe
    writePipe.Close()

    parent.Wait()
    os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
    command := strings.Join(cmdArray, "")
    log.Infof("Full command: %s", command)
    writePipe.WriteString(command)
}

