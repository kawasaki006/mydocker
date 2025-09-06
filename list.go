package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path"
    "text/tabwriter"

    "github.com/kawasaki006/mydocker/container"
     log "github.com/sirupsen/logrus"
)

func listContainers() {
    // read all container config files
    files, err := os.ReadDir(container.InfoLocation)
    if err != nil {
        log.Errorf("error reading info directory [%s]: %v", container.InfoLocation, err)
    }
    containers := make([]*container.ContainerInfo, 0, len(files))
    for _, file := range files {
        tmp, err := fetchInfo(file)
        if err != nil {
            log.Errorf("error fetching info: %v", err)
            continue
        }
        containers = append(containers, tmp)
    }
    //print
    w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
    _, err = fmt.Fprint(w, "ID\tNAME\tPID\tIP\tSTATUS\tCOMMAND\tCREATED\n")
    if err != nil {
        log.Errorf("Fprint containers error: %v", err)
    }
    for _, i := range containers {
        _,err = fmt.Fprint(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
            i.Id,
            i.Name,
            i.Pid,
            i.Status,
            i.Command,
            i.CreatedTime)
        if err != nil {
            log.Errorf("Fprint error: %v", err)
        }
    }
    if err = w.Flush(); err != nil {
        log.Errorf("error flushing: %v", err)
    }
}

func fetchInfo(file os.DirEntry) (*container.ContainerInfo, error) {
    configFileDir := path.Join(container.InfoLocation, file.Name())
    configFileDir = path.Join(configFileDir, container.ConfigName)
    // read file
    infoJson, err := os.ReadFile(configFileDir)
    if err != nil {
        log.Errorf("error reading file [%s]: %v", configFileDir, err)
        return nil, err
    }
    info := new(container.ContainerInfo)
    // decode json
    if err = json.Unmarshal(infoJson, info); err != nil {
        log.Errorf("error decoding json: %v", err)
        return nil, err
    }

    return info, nil
}
