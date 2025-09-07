package container

import (
    "fmt"
    "os"
    "path"
    "math/rand"
    "strings"
    "strconv"
    "time"
    "encoding/json"
)

const (
    RUNNING      = "running"
    STOP         = "stopped"
    EXIT         = "exited"
    InfoLocation = "/var/lib/mydocker/containers/"
    ConfigName   = "config.json"
    IdLength     = 10
    LogFile      = "%s-json.log"
)

type ContainerInfo struct {
    Pid         string    `json:"pid"`
    Id          string    `json:"id"`
    Name        string    `json:"name"`
    Command     string    `json:"command"`
    CreatedTime string    `json:"createdTime"`
    Status      string    `json:"status"`
    Volume      string    `json:"volume"`
}

func RecordContainerInfo(containerPid int, commandArray []string, containerName, containerId, volume string) (*ContainerInfo, error) {
    if containerName == "" {
        containerName = containerId
    }
    command := strings.Join(commandArray, "")
    containerInfo := &ContainerInfo{
        Pid:         strconv.Itoa(containerPid),
        Id:          containerId,
        Name:        containerName,
        Command:     command,
        CreatedTime: time.Now().Format("2025-09-05 12:00:00"),
        Status:      RUNNING,
        Volume:      volume,
    }

    // encode into json
    infoJsonBytes, err := json.Marshal(containerInfo)
    if err != nil {
        return containerInfo, fmt.Errorf("error encoding info object into json: %v", err)
    }
    infoJsonStr := string(infoJsonBytes)

    // mkdir for container info
    dirPath := path.Join(InfoLocation, containerId)
    if _, err := os.Stat(dirPath); err != nil {
        // dir not exist, make it
        if err := os.Mkdir(dirPath, 0622); err != nil {
            return containerInfo, fmt.Errorf("error creating container info dir at [%s]: %v", dirPath, err)
        }
    }
    // write file to path
    filePath := path.Join(dirPath, ConfigName)
    file, err := os.Create(filePath)
    if err != nil {
        return containerInfo, fmt.Errorf("error creating config file: %v", err)
    }
    defer file.Close()
    if _, err = file.WriteString(infoJsonStr); err != nil {
        return containerInfo, fmt.Errorf("error writing to config file: %v", err)
    }
    return containerInfo, nil
}

func DeleteContainerInfo(containerId string) error {
    dirPath := path.Join(InfoLocation, containerId)
    if err := os.RemoveAll(dirPath); err != nil {
        return fmt.Errorf("error cleaning up container info config: %v", err)
    }
    return nil
}

func GenerateContainerId() string {
    return randStringBytes(IdLength)
}

func randStringBytes(n int) string {
    letterBytes := "1234567890"
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func GetLogFile(containerId string) string {
    return fmt.Sprintf(LogFile, containerId)
}
