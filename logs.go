package main

import (
    "path"
    "os"
    "io"
    "fmt"

    "github.com/kawasaki006/mydocker/container"
    log "github.com/sirupsen/logrus"
)

func logContainer(containerId string) {
    logFilePath := path.Join(container.InfoLocation, containerId) + container.GetLogFile(containerId)
    file, err := os.Open(logFilePath)
    defer file.Close()
    if err != nil {
        log.Errorf("error opening log file [%s]: %v", logFilePath, err)
        return
    }
    content, err := io.ReadAll(file)
    if err != nil {
        log.Errorf("error reading log file [%s]: %v", logFilePath, err)
        return
    }
    _, err = fmt.Fprint(os.Stdout, string(content))
    if err != nil {
        log.Errorf("error printing log: %v", err)
        return
    }
}
