package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path"
    "strings"

    "github.com/kawasaki006/mydocker/container"
    _ "github.com/kawasaki006/mydocker/nsenter"
    log "github.com/sirupsen/logrus"
)

const (
    EnvExecPid = "mydocker_pid"
    EnvExecCmd = "mydocker_cmd"
)

func ExecContainer(containerId string, comArray []string) {
    pid, err := getPidByContainerId(containerId)
    if err != nil {
        log.Errorf("Exec container getContainerPidByName %s error %v", containerId, err)
        return
    }
    // create a new proc
    cmd := exec.Command("/proc/self/exe", "exec")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    // set env pid and command
    cmdStr := strings.Join(comArray, " ")
    log.Infof("container pid: %s, command: %s", pid, cmdStr)
    _ = os.Setenv(EnvExecPid, pid)
    _ = os.Setenv(EnvExecCmd, cmdStr)
    // copy and set container's envs
    containerEnvs := getEnvsByPid(pid)
    cmd.Env = append(os.Environ(), containerEnvs...)
    // run the proc
    if err = cmd.Run(); err != nil {
        log.Errorf("error exec container %s: %v", containerId, err)
    }
}

func getPidByContainerId(containerId string) (string, error) {
    // get container config
    dirPath := path.Join(container.InfoLocation, containerId)
    configFilePath := path.Join(dirPath, container.ConfigName)
    // read config
    content, err := os.ReadFile(configFilePath)
    if err != nil {
        return "", fmt.Errorf("error reading config file: %v", err)
    }
    var containerInfo container.ContainerInfo
    if err = json.Unmarshal(content, &containerInfo); err != nil {
        return "", fmt.Errorf("error decoding config file: %v", err)
    }
    return containerInfo.Pid, nil
}

func getEnvsByPid(pid string) []string {
    path := fmt.Sprintf("/proc/%s/environ", pid)
    content,err := os.ReadFile(path)
    if err != nil {
        log.Errorf("error reading envs [%s]: %v", path, err)
        return nil
    }
    envs := strings.Split(string(content), "\u0000")
    return envs
}
